/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openapi

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/routes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/pkg/validation/spec"

	apiextensionshelpers "k8s.io/apiextensions-apiserver/pkg/apihelpers"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	informers "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions/apiextensions/v1"
	listers "k8s.io/apiextensions-apiserver/pkg/client/listers/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/controller/openapi/builder"
)

// Controller watches CustomResourceDefinitions and publishes validation schema
type Controller struct {
	crdLister  listers.CustomResourceDefinitionLister
	crdIndexer cache.Indexer
	crdsSynced cache.InformerSynced

	// To allow injection for testing.
	syncFn func(string) error

	queue workqueue.RateLimitingInterface

	staticSpec             *spec.Swagger
	openAPIServiceProvider routes.OpenAPIServiceProvider

	// specs per scope and per version and per CRD name
	lock     sync.Mutex
	crdSpecs map[string]map[string]map[string]*spec.Swagger
}

// NewController creates a new Controller with input CustomResourceDefinition informer
func NewController(crdInformer informers.CustomResourceDefinitionInformer) *Controller {
	c := &Controller{
		crdLister:  crdInformer.Lister(),
		crdIndexer: crdInformer.Informer().GetIndexer(),
		crdsSynced: crdInformer.Informer().HasSynced,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "crd_openapi_controller"),
		crdSpecs:   map[string]map[string]map[string]*spec.Swagger{},
	}

	crdInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addCustomResourceDefinition,
		UpdateFunc: c.updateCustomResourceDefinition,
		DeleteFunc: c.deleteCustomResourceDefinition,
	})

	c.syncFn = c.sync
	return c
}

// HACK:
//  Everything regarding OpenAPI and resource discovery is managed through controllers currently
// (a number of controllers highly coupled with the corresponding http handlers).
// The following code is an attempt at provising CRD tenancy while accommodating the current design without being too much invasive,
// because doing differently would have meant too much refactoring..
// But in the long run the "do this dynamically, not as part of a controller" is probably going to be important.
// openapi/crd generation is expensive, so doing on a controller means that CPU and memory scale O(crds),
// when we really want them to scale O(active_scopes).

func (c *Controller) setScopeCrdSpecs(scopeName, crdName string, newSpecs map[string]*spec.Swagger) {
	_, found := c.crdSpecs[scopeName]
	if !found {
		c.crdSpecs[scopeName] = map[string]map[string]*spec.Swagger{}
	}
	c.crdSpecs[scopeName][crdName] = newSpecs
	c.openAPIServiceProvider.AddScope(scopeName)
}

func (c *Controller) removeScopeCrdSpecs(scope rest.Scope, crdName string) bool {
	scopeName := scope.Name()
	if _, found := c.crdSpecs[scopeName]; !found {
		return false
	}
	if _, found := c.crdSpecs[scopeName][crdName]; !found {
		return false
	}

	delete(c.crdSpecs[scopeName], crdName)

	if len(c.crdSpecs[scopeName]) == 0 {
		delete(c.crdSpecs, scopeName)
		c.openAPIServiceProvider.RemoveScope(scopeName)
	}

	return true
}

func (c *Controller) getClusterCrdSpecs(scope rest.Scope, crdName string) (map[string]*spec.Swagger, bool) {
	scopeName := scope.Name()
	_, specsFoundForCluster := c.crdSpecs[scopeName]
	if !specsFoundForCluster {
		return map[string]*spec.Swagger{}, false
	}
	crdSpecs, found := c.crdSpecs[scopeName][crdName]
	return crdSpecs, found
}

// Run sets openAPIAggregationManager and starts workers
func (c *Controller) Run(staticSpec *spec.Swagger, openAPIServiceProvider routes.OpenAPIServiceProvider, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()
	defer klog.Infof("Shutting down OpenAPI controller")

	klog.Infof("Starting OpenAPI controller")

	c.staticSpec = staticSpec
	c.openAPIServiceProvider = openAPIServiceProvider

	if !cache.WaitForCacheSync(stopCh, c.crdsSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	scopes := c.crdIndexer.ListIndexFuncValues(cache.ListAllIndex)
	for _, scopeName := range scopes {
		scope := cache.NewScope(scopeName)
		// create initial spec to avoid merging once per CRD on startup
		crds, err := c.crdLister.Scoped(scope).List(labels.Everything())
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("failed to initially list all CRDs: %v", err))
			return
		}
		for _, crd := range crds {
			if !apiextensionshelpers.IsCRDConditionTrue(crd, apiextensionsv1.Established) {
				continue
			}

			newSpecs, changed, err := buildVersionSpecs(crd, nil)
			if err != nil {
				utilruntime.HandleError(fmt.Errorf("failed to build OpenAPI spec of CRD %s: %v", crd.Name, err))
			} else if !changed {
				continue
			}

			c.setScopeCrdSpecs(scopeName, crd.Name, newSpecs)
		}
		if err := c.updateSpecLocked(); err != nil {
			utilruntime.HandleError(fmt.Errorf("failed to initially create OpenAPI spec for CRDs: %v", err))
			return
		}
	}

	// only start one worker thread since its a slow moving API
	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	// log slow aggregations
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		if elapsed > time.Second {
			klog.Warningf("slow openapi aggregation of %q: %s", key.(string), elapsed)
		}
	}()

	err := c.syncFn(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("%v failed with: %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

func (c *Controller) sync(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	scope, err := cache.ScopeFromKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("error getting scope from key %q: %w", key, err))
		return nil
	}

	crd, err := c.crdLister.Get(key)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	// do we have to remove all specs of this CRD?
	if errors.IsNotFound(err) || !apiextensionshelpers.IsCRDConditionTrue(crd, apiextensionsv1.Established) {
		var crdName string
		if crd != nil {
			crdName = crd.Name
		} else {
			queueKey, err := cache.DecodeKeyFunc(key)
			if err != nil {
				utilruntime.HandleError(fmt.Errorf("error decoding cache key %q: %w", key, err))
				return nil
			}
			crdName = queueKey.Name()
		}
		if !c.removeScopeCrdSpecs(scope, crdName) {
			return nil
		}
		klog.V(2).Infof("Updating CRD OpenAPI spec because %s was removed", key)
		regenerationCounter.With(map[string]string{"crd": key, "reason": "remove"})
		return c.updateSpecLocked()
	}

	// compute CRD spec and see whether it changed
	oldSpecs, updated := c.getClusterCrdSpecs(scope, crd.Name)
	newSpecs, changed, err := buildVersionSpecs(crd, oldSpecs)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}

	// update specs of this CRD
	c.setScopeCrdSpecs(scope.Name(), crd.Name, newSpecs)
	klog.V(2).Infof("Updating CRD OpenAPI spec because %s changed", key)
	reason := "add"
	if updated {
		reason = "update"
	}
	regenerationCounter.With(map[string]string{"crd": key, "reason": reason})
	return c.updateSpecLocked()
}

func buildVersionSpecs(crd *apiextensionsv1.CustomResourceDefinition, oldSpecs map[string]*spec.Swagger) (map[string]*spec.Swagger, bool, error) {
	newSpecs := map[string]*spec.Swagger{}
	anyChanged := false
	for _, v := range crd.Spec.Versions {
		if !v.Served {
			continue
		}
		// Defaults are not pruned here, but before being served.
		spec, err := builder.BuildOpenAPIV2(crd, v.Name, builder.Options{V2: true})
		if err != nil {
			return nil, false, err
		}
		newSpecs[v.Name] = spec
		if oldSpecs[v.Name] == nil || !reflect.DeepEqual(oldSpecs[v.Name], spec) {
			anyChanged = true
		}
	}
	if !anyChanged && len(oldSpecs) == len(newSpecs) {
		return newSpecs, false, nil
	}

	return newSpecs, true, nil
}

// updateSpecLocked aggregates all OpenAPI specs and updates openAPIService.
// It is not thread-safe. The caller is responsible to hold proper lock (Controller.lock).
func (c *Controller) updateSpecLocked() error {
	var errs []error
	for scopeName, scopeSpecs := range c.crdSpecs {
		crdSpecs := []*spec.Swagger{}
		for _, versionSpecs := range scopeSpecs {
			for _, s := range versionSpecs {
				crdSpecs = append(crdSpecs, s)
			}
		}
		mergedSpec, err := builder.MergeSpecs(c.staticSpec, crdSpecs...)
		if err != nil {
			return fmt.Errorf("failed to merge specs: %w", err)
		}
		if err := c.openAPIServiceProvider.ForScope(scopeName).UpdateSpec(mergedSpec); err != nil {
			errs = append(errs, err)
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (c *Controller) addCustomResourceDefinition(obj interface{}) {
	castObj := obj.(*apiextensionsv1.CustomResourceDefinition)
	klog.V(4).Infof("Adding customresourcedefinition %s", castObj.Name)
	c.enqueue(castObj)
}

func (c *Controller) updateCustomResourceDefinition(oldObj, newObj interface{}) {
	castNewObj := newObj.(*apiextensionsv1.CustomResourceDefinition)
	klog.V(4).Infof("Updating customresourcedefinition %s", castNewObj.Name)
	c.enqueue(castNewObj)
}

func (c *Controller) deleteCustomResourceDefinition(obj interface{}) {
	castObj, ok := obj.(*apiextensionsv1.CustomResourceDefinition)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			klog.Errorf("Couldn't get object from tombstone %#v", obj)
			return
		}
		castObj, ok = tombstone.Obj.(*apiextensionsv1.CustomResourceDefinition)
		if !ok {
			klog.Errorf("Tombstone contained object that is not expected %#v", obj)
			return
		}
	}
	klog.V(4).Infof("Deleting customresourcedefinition %q", castObj.Name)
	c.enqueue(castObj)
}

func (c *Controller) enqueue(obj *apiextensionsv1.CustomResourceDefinition) {
	key, _ := cache.ObjectKeyFunc(obj)
	c.queue.Add(key)
}
