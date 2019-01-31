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
	"time"

	"github.com/go-openapi/spec"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	informers "k8s.io/apiextensions-apiserver/pkg/client/informers/internalversion/apiextensions/internalversion"
	listers "k8s.io/apiextensions-apiserver/pkg/client/listers/apiextensions/internalversion"
	apiextensionsfeatures "k8s.io/apiextensions-apiserver/pkg/features"
)

// AggregationManager is the interface between OpenAPI Aggregator service and a controller
// that manages CRD openapi spec aggregation
type AggregationManager interface {
	// AddUpdateLocalAPIService allows adding/updating local API service with nil handler and
	// nil Spec.Service. This function can be used for local dynamic OpenAPI spec aggregation
	// management (e.g. CRD)
	AddUpdateLocalAPIServiceSpec(name string, spec *spec.Swagger, etag string) error
	RemoveAPIServiceSpec(apiServiceName string) error
}

// Controller watches CustomResourceDefinitions and publishes validation schema
type Controller struct {
	crdLister  listers.CustomResourceDefinitionLister
	crdsSynced cache.InformerSynced

	// To allow injection for testing.
	syncFn func(gvk schema.GroupVersionKind) error

	queue workqueue.RateLimitingInterface

	openAPIAggregationManager AggregationManager
}

// NewController creates a new Controller with input CustomResourceDefinition informer
func NewController(crdInformer informers.CustomResourceDefinitionInformer) *Controller {
	c := &Controller{
		crdLister:  crdInformer.Lister(),
		crdsSynced: crdInformer.Informer().HasSynced,

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "crd_openapi_controller"),
	}

	crdInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addCustomResourceDefinition,
		UpdateFunc: c.updateCustomResourceDefinition,
		DeleteFunc: c.deleteCustomResourceDefinition,
	})

	c.syncFn = c.sync
	return c
}

// Run sets openAPIAggregationManager and starts workers
func (c *Controller) Run(stopCh <-chan struct{}, aggregationManager AggregationManager) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()
	defer klog.Infof("Shutting down CRDOpenAPIController")

	klog.Infof("Starting CRDOpenAPIController")

	c.openAPIAggregationManager = aggregationManager

	if !cache.WaitForCacheSync(stopCh, c.crdsSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
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

	err := c.syncFn(key.(schema.GroupVersionKind))
	if err == nil {
		c.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("%v failed with: %v", key, err))
	c.queue.AddRateLimited(key)
	return true
}

func (c *Controller) sync(gvk schema.GroupVersionKind) error {
	if c.openAPIAggregationManager == nil || !utilfeature.DefaultFeatureGate.Enabled(apiextensionsfeatures.CustomResourceValidation) {
		return nil
	}

	crds, err := c.crdLister.List(labels.Everything())
	if err != nil {
		return err
	}
	specName := gvk.String()
	found := false

	for _, crd := range crds {
		if !apiextensions.IsCRDConditionTrue(crd, apiextensions.Established) {
			continue
		}
		if crd.Spec.Group != gvk.Group || crd.Spec.Names.Kind != gvk.Kind {
			continue
		}
		if !apiextensions.HasVersionServed(crd, gvk.Version) {
			continue
		}

		found = true

		start := time.Now()
		spec, etag, err := BuildSwagger(crd, gvk.Version)
		if err != nil {
			return err
		}
		if err := c.openAPIAggregationManager.AddUpdateLocalAPIServiceSpec(specName, spec, etag); err != nil {
			return err
		}

		elapsed := time.Since(start)
		klog.Errorf(">>>>> in total, build and aggregate openapi took %s", elapsed)
	}

	if !found {
		return c.openAPIAggregationManager.RemoveAPIServiceSpec(specName)
	}
	return nil
}

func (c *Controller) addCustomResourceDefinition(obj interface{}) {
	castObj := obj.(*apiextensions.CustomResourceDefinition)
	klog.V(4).Infof("Adding customresourcedefinition %s", castObj.Name)
	c.enqueue(castObj)
}

func (c *Controller) updateCustomResourceDefinition(oldObj, newObj interface{}) {
	castNewObj := newObj.(*apiextensions.CustomResourceDefinition)
	castOldObj := oldObj.(*apiextensions.CustomResourceDefinition)
	klog.V(4).Infof("Updating customresourcedefinition %s", castOldObj.Name)
	c.enqueue(castNewObj)
	c.enqueue(castOldObj)
}

func (c *Controller) deleteCustomResourceDefinition(obj interface{}) {
	castObj, ok := obj.(*apiextensions.CustomResourceDefinition)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			klog.Errorf("Couldn't get object from tombstone %#v", obj)
			return
		}
		castObj, ok = tombstone.Obj.(*apiextensions.CustomResourceDefinition)
		if !ok {
			klog.Errorf("Tombstone contained object that is not expected %#v", obj)
			return
		}
	}
	klog.V(4).Infof("Deleting customresourcedefinition %q", castObj.Name)
	c.enqueue(castObj)
}

func (c *Controller) enqueue(obj *apiextensions.CustomResourceDefinition) {
	for _, v := range obj.Spec.Versions {
		c.queue.Add(schema.GroupVersionKind{Group: obj.Spec.Group, Version: v.Name, Kind: obj.Spec.Names.Kind})
	}
}
