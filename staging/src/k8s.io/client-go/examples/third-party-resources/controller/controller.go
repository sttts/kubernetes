/*
Copyright 2017 The Kubernetes Authors.

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

package controller

import (
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	tprv1 "k8s.io/client-go/examples/third-party-resources/apis/tpr/v1"
)

// Watcher is an example of watching on resource create/update/delete events
type ExampleController struct {
	Examples cache.Store

	exampleClient *rest.RESTClient
	informer      cache.Controller
}

func NewExampleController(exampleClient *rest.RESTClient) (*ExampleController, error) {
	c := &ExampleController{
		exampleClient: exampleClient,
	}

	lw := cache.NewListWatchFromClient(
		c.exampleClient,
		tprv1.ExampleResourcePlural,
		apiv1.NamespaceAll,
		fields.Everything())

	c.Examples, c.informer = cache.NewInformer(
		lw,
		&tprv1.Example{}, // The object type

		// resyncPeriod
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,

		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		},
	)

	return c, nil
}

// Run starts an Example resource watcher
func (c *ExampleController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()

	fmt.Print("Watch Example objects\n")

	go c.informer.Run(stopCh)

	<-stopCh
}

func (c *ExampleController) onAdd(obj interface{}) {
	example := obj.(*tprv1.Example)
	fmt.Printf("[WATCH] OnAdd %s\n", example.SelfLink)
}

func (c *ExampleController) onUpdate(oldObj, newObj interface{}) {
	oldExample := oldObj.(*tprv1.Example)
	newExample := newObj.(*tprv1.Example)
	fmt.Printf("[WATCH] OnUpdate oldObj: %s\n", oldExample.SelfLink)
	fmt.Printf("[WATCH] OnUpdate newObj: %s\n", newExample.SelfLink)
}

func (c *ExampleController) onDelete(obj interface{}) {
	example := obj.(*tprv1.Example)
	fmt.Printf("[WATCH] OnDelete %s\n", example.SelfLink)
}
