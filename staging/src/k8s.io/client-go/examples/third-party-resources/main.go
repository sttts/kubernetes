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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	exampletprv1 "k8s.io/client-go/examples/third-party-resources/apis/tpr/v1"
	exampleclient "k8s.io/client-go/examples/third-party-resources/client"
	examplecontroller "k8s.io/client-go/examples/third-party-resources/controller"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := buildConfig(*kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// initialize third party resource if it does not exist
	err = exampleclient.CreateTPR(clientset)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}

	// make a new config for our extension's API group, using the first config as a baseline
	exampleClient, _, err := exampleclient.NewClient(config)
	if err != nil {
		panic(err)
	}

	// wait until TPR gets processed
	exampleclient.WaitForExampleResource(exampleClient)

	// start a controller on instances of our TPR
	controller, err := examplecontroller.NewExampleController(exampleClient)
	if err != nil {
		panic(err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	go controller.Run(stopCh)

	// Create an instance of our TPR
	example := &exampletprv1.Example{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example1",
		},
		Spec: exampletprv1.ExampleSpec{
			Foo: "hello",
			Bar: true,
		},
	}
	var result exampletprv1.Example
	err = exampleClient.Post().
		Resource(exampletprv1.ExampleResourcePlural).
		Namespace(apiv1.NamespaceDefault).
		Body(example).
		Do().Into(&result)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v\n", result)
	} else {
		panic(err)
	}

	// Fetch a list of our TPRs
	exampleList := exampletprv1.ExampleList{}
	err = exampleClient.Get().Resource(exampletprv1.ExampleResourcePlural).Do().Into(&exampleList)
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleList)

	// the created TPR should show up in the controller store
	var exampleFromStore *exampletprv1.Example
	err = wait.Poll(100*time.Millisecond, 30*time.Second, func() (bool, error) {
		obj, exists, err := controller.Examples.GetByKey(apiv1.NamespaceDefault + "/" + example.Name)
		if !exists || err != nil {
			return exists, err
		}
		exampleFromStore = obj.(*exampletprv1.Example)
		return true, nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("FOUND IN STORE: %#v\n", exampleFromStore)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
