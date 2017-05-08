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

package thirdparty

import (
	"reflect"
	"testing"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	corev1 "k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/test/integration/framework"

	exampletprv1 "k8s.io/client-go/examples/third-party-resources/apis/tpr/v1"
	exampleclient "k8s.io/client-go/examples/third-party-resources/client"
	examplecontroller "k8s.io/client-go/examples/third-party-resources/controller"
)

func TestClientGoThirdPartyResourceExample(t *testing.T) {
	_, s := framework.RunAMaster(framework.NewIntegrationTestMasterConfig())
	defer s.Close()

	config := &rest.Config{Host: s.URL, ContentConfig: rest.ContentConfig{NegotiatedSerializer: api.Codecs}}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("Creating TPR %q", exampletprv1.ExampleResourcePlural)
	if err := exampleclient.CreateTPR(clientset); err != nil {
		t.Fatalf("unexpected error creating the ThirdPartyResource: %v", err)
	}

	exampleClient, _, err := exampleclient.NewClient(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("Waiting for TPR %q to show up", exampletprv1.ExampleResourcePlural)
	if err := exampleclient.WaitForExampleResource(exampleClient); err != nil {
		t.Fatalf("ThirdPartyResource examples did not show up: %v", err)
	}
	t.Logf("TPR %q is active", exampletprv1.ExampleResourcePlural)

	t.Logf("Starting a controller on instances of TPR %q", exampletprv1.ExampleResourcePlural)
	controller, err := examplecontroller.NewExampleController(exampleClient)
	if err != nil {
		t.Fatalf("Failed to launch ExampleController: %v", err)
	}
	stopCh := make(chan struct{})
	defer close(stopCh)
	go controller.Run(stopCh)

	t.Logf("Creating example instance")
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
		Namespace(corev1.NamespaceDefault).
		Body(example).
		Do().Into(&result)
	if err != nil && apierrors.IsAlreadyExists(err) {
		t.Fatalf("Failed to create TPR object: %v", err)
	}

	// Fetch a list of our TPRs
	t.Logf("Checking that the example instance shows up in a LIST request")
	exampleList := exampletprv1.ExampleList{}
	err = exampleClient.Get().Resource(exampletprv1.ExampleResourcePlural).Do().Into(&exampleList)
	if err != nil {
		t.Fatalf("Failed to fetch a list of examples: %v", err)
	}
	if len(exampleList.Items) != 1 {
		t.Fatalf("Expected exactly one example in list, got: %#v", exampleList)
	}
	if !reflect.DeepEqual(example.Spec, exampleList.Items[0].Spec) {
		t.Fatalf("Didn't find example with the original spec: %v", diff.ObjectDiff(example, exampleList.Items[0].Spec))
	}

	// the created TPR should show up in the controller store
	t.Logf("Waiting for example instance to show up in store")
	var exampleFromStore *exampletprv1.Example
	err = wait.Poll(100*time.Millisecond, 30*time.Second, func() (bool, error) {
		obj, exists, err := controller.Examples.GetByKey(corev1.NamespaceDefault + "/" + example.Name)
		if !exists || err != nil {
			return exists, err
		}
		exampleFromStore = obj.(*exampletprv1.Example)
		return true, nil
	})
	if err != nil {
		t.Fatalf("example did not show up in store: %v", err)
	}

	t.Logf("Found example in store: %#v", exampleFromStore)
}
