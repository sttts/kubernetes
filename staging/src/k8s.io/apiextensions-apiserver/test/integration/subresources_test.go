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

package integration

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilfeature "k8s.io/apiserver/pkg/util/feature"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/test/integration/testserver"
)

func newNoxuSubResourcesCRD(scope apiextensionsv1beta1.ResourceScope) *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: "noxus.mygroup.example.com"},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   "mygroup.example.com",
			Version: "v1beta1",
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural:     "noxus",
				Singular:   "nonenglishnoxu",
				Kind:       "WishIHadChosenNoxu",
				ShortNames: []string{"foo", "bar", "abc", "def"},
				ListKind:   "NoxuItemList",
			},
			Scope: apiextensionsv1beta1.NamespaceScoped,
			SubResources: &apiextensionsv1beta1.CustomResourceSubResources{
				Status: &apiextensionsv1beta1.CustomResourceSubResourceStatus{},
				Scale: &apiextensionsv1beta1.CustomResourceSubResourceScale{
					SpecReplicasPath:   ".spec.replicas",
					StatusReplicasPath: ".status.replicas",
					LabelSelectorPath:  ".spec.labelSelector",
					ScaleGroupVersion:  "autoscaling/v1",
				},
			},
		},
	}
}

func newNoxuSubResourceInstance(namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "mygroup.example.com/v1beta1",
			"kind":       "WishIHadChosenNoxu",
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
			"spec": map[string]interface{}{
				"num":      int64(10),
				"replicas": int64(3),
			},
		},
	}
}

func TestStatusSubResource(t *testing.T) {
	// enable alpha feature CustomResourceSubResources
	if err := utilfeature.DefaultFeatureGate.Set("CustomResourceSubResources=true"); err != nil {
		t.Errorf("failed to enable feature gate for CustomResourceSubResources: %v", err)
	}

	stopCh, apiExtensionClient, clientPool, err := testserver.StartDefaultServer()
	if err != nil {
		t.Fatal(err)
	}
	defer close(stopCh)

	noxuDefinition := newNoxuSubResourcesCRD(apiextensionsv1beta1.NamespaceScoped)
	noxuVersionClient, err := testserver.CreateNewCustomResourceDefinition(noxuDefinition, apiExtensionClient, clientPool)
	if err != nil {
		t.Fatal(err)
	}

	ns := "not-the-default"
	noxuResourceClient := NewNamespacedCustomResourceClient(ns, noxuVersionClient, noxuDefinition)
	_, err = instantiateCustomResource(t, newNoxuSubResourceInstance(ns, "foo"), noxuResourceClient, noxuDefinition)
	if err != nil {
		t.Fatalf("unable to create noxu instance: %v", err)
	}

	gottenNoxuInstance, err := noxuResourceClient.Get("foo", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	// .status.num = 10
	ok := unstructured.SetNestedField(gottenNoxuInstance.Object, int64(10), "status", "num")
	if !ok {
		t.Fatalf("unable to set status")
	}

	// .spec.num = 20
	ok = unstructured.SetNestedField(gottenNoxuInstance.Object, int64(20), "spec", "num")
	if !ok {
		t.Fatalf("unable to set spec")
	}

	updatedStatusInstance, err := noxuResourceClient.UpdateStatus(gottenNoxuInstance)
	if err != nil {
		t.Fatalf("unable to update status: %v", err)
	}

	// UpdateStatus should not update spec. Check that .spec.num remains 10.
	specNum, ok := unstructured.NestedInt64(updatedStatusInstance.Object, "spec", "num")
	if !ok {
		t.Fatalf("unable to get .spec.num")
	}
	if specNum != int64(10) {
		t.Fatalf("expected %v, got %v", int64(10), specNum)
	}

	statusNum, ok := unstructured.NestedInt64(updatedStatusInstance.Object, "status", "num")
	if !ok {
		t.Fatalf("unable to get .status.num")
	}
	if statusNum != int64(10) {
		t.Fatalf("expected %v, got %v", int64(10), statusNum)
	}
}

func TestScaleSubResource(t *testing.T) {
	// enable alpha feature CustomResourceSubResources
	if err := utilfeature.DefaultFeatureGate.Set("CustomResourceSubResources=true"); err != nil {
		t.Errorf("failed to enable feature gate for CustomResourceSubResources: %v", err)
	}

	stopCh, apiExtensionClient, clientPool, err := testserver.StartDefaultServer()
	if err != nil {
		t.Fatal(err)
	}
	defer close(stopCh)

	noxuDefinition := newNoxuSubResourcesCRD(apiextensionsv1beta1.NamespaceScoped)
	noxuVersionClient, err := testserver.CreateNewCustomResourceDefinition(noxuDefinition, apiExtensionClient, clientPool)
	if err != nil {
		t.Fatal(err)
	}

	ns := "not-the-default"
	noxuResourceClient := NewNamespacedCustomResourceClient(ns, noxuVersionClient, noxuDefinition)
	_, err = instantiateCustomResource(t, newNoxuSubResourceInstance(ns, "foo"), noxuResourceClient, noxuDefinition)
	if err != nil {
		t.Fatalf("unable to create noxu instance: %v", err)
	}

	// TODO: fix this
	_, err = noxuResourceClient.GetScale("foo", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("unable to get scale: %v", err)
	}
}
