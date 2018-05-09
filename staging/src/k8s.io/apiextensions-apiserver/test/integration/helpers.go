/*
Copyright 2018 The Kubernetes Authors.

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
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// UpdateCustomResourceDefinition updates a CRD, retrying up to 5 times on version conflict errors.
func UpdateCustomResourceDefinition(client clientset.Interface, name string, update func(*apiextensionsv1beta1.CustomResourceDefinition)) (*apiextensionsv1beta1.CustomResourceDefinition, error) {
	for i := 0; i < 5; i++ {
		crd, err := client.ApiextensionsV1beta1().CustomResourceDefinitions().Get(name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to get Service %q: %v", name, err)
		}
		update(crd)
		crd, err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Update(crd)
		if err == nil {
			return crd, nil
		}
		if !errors.IsConflict(err) && !errors.IsServerTimeout(err) {
			return nil, fmt.Errorf("failed to update CustomResourceDefinition %q: %v", name, err)
		}
	}
	return nil, fmt.Errorf("too many retries updating CustomResourceDefinition %q", name)
}
