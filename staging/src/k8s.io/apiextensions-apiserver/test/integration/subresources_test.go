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
				"num": 10,
			},
		},
	}
}

func TestStatusSubResource(t *testing.T) {
	// enable alpha feature CustomResourceDefaulting
	if err := utilfeature.DefaultFeatureGate.Set("CustomResourceSubResources=true"); err != nil {
		t.Errorf("failed to enable feature gate for CustomResourceDefaulting: %v", err)
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

	gottenNoxuInstance.Object = map[string]interface{}{
		"apiVersion": "mygroup.example.com/v1beta1",
		"kind":       "WishIHadChosenNoxu",
		"metadata": map[string]interface{}{
			"namespace": "not-the-default",
			"name":      "foo",
		},
		"spec": map[string]interface{}{
			"num": 10,
		},
		"status": map[string]interface{}{
			"num": 10,
		},
	}

	// TODO: fix this error
	_, err = noxuResourceClient.UpdateStatus(gottenNoxuInstance)
	if err != nil {
		t.Fatalf("unable to update status: %v", err)
	}
}
