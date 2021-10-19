/*
Copyright 2020 The Kubernetes Authors.

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

package apimachinery

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/onsi/ginkgo"
)

var _ = SIGDescribe("KCP MultiClusterWorkflow", func() {
	f := framework.NewDefaultFramework("controller-workflow")

	/*
	   Release:
	   Testname: Operate on objects from multi-cluster list
	   Description: Ensure that an object returned from the multi-cluster
	   list can be used in mutating requests in a single-cluster context.
	*/
	framework.ConformanceIt("should use an object from multi-cluster list in single-cluster contexts", func() {
		ns := f.Namespace.Name
		ginkgo.By("creating configmaps in single-cluster contexts")
		for i, c := range []clientset.Interface{f.ClientSet, f.SecondaryClientSet, f.TertiaryClientSet} {
			ginkgo.By(fmt.Sprintf("running with clientset %s", []string{"primary", "secondary", "tertiary"}[i]))
			if c != nil {
				for i := 0; i < 5; i++ {
					testConfigMap := &v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name: fmt.Sprintf("e2e-workflow-test-configmap-%d", i),
							Labels: map[string]string{
								watchConfigMapNamespaceLabelKey: ns,
							},
						},
					}

					_, err := c.CoreV1().ConfigMaps(ns).Create(context.TODO(), testConfigMap, metav1.CreateOptions{})
					framework.ExpectNoError(err, "failed to create a configmap in namespace: %s", ns)
				}
			}
		}

		ginkgo.By("fetching a list of the configmaps from a cross-cluster context")
		configMaps, err := f.MultiClusterClientSet.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labels.Set{watchConfigMapNamespaceLabelKey: ns}).String(),
		})
		framework.ExpectNoError(err, "listing configmaps in all namespaces")

		ginkgo.By("updating configmaps in single-cluster contexts")
		for _, configMap := range configMaps.Items {
			toUpdate := configMap.DeepCopy()
			toUpdate.Labels["mutated"] = "byUpdate"
			_, err = f.ClusterlessClientSet.Cluster(configMap.ClusterName).CoreV1().ConfigMaps(ns).Update(context.TODO(), toUpdate, metav1.UpdateOptions{}) // preconditions are implicit
			framework.ExpectNoError(err, "updating configmap %q in namespace %q for cluster %q", toUpdate.Name, toUpdate.Namespace, toUpdate.ClusterName)

			oldData, err := json.Marshal(toUpdate)
			framework.ExpectNoError(err, "failed to marshal unmodified configmap %q in namespace %q for cluster %q into JSON", toUpdate.Name, toUpdate.Namespace, toUpdate.ClusterName)

			toUpdate.Labels["mutated"] = "byPatch"

			newData, err := json.Marshal(toUpdate)
			framework.ExpectNoError(err, "failed to marshal unmodified configmap %q in namespace %q for cluster %q into JSON", toUpdate.Name, toUpdate.Namespace, toUpdate.ClusterName)

			patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.ConfigMap{})
			framework.ExpectNoError(err, "failed to create two way merge patch for configmap %q in namespace %q for cluster %q", toUpdate.Name, toUpdate.Namespace, toUpdate.ClusterName)

			_, err = f.ClusterlessClientSet.Cluster(configMap.ClusterName).CoreV1().ConfigMaps(ns).Patch(context.TODO(), toUpdate.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{}) // preconditions are not supported
			framework.ExpectNoError(err, "patching configmap %q in namespace %q for cluster %q", toUpdate.Name, toUpdate.Namespace, toUpdate.ClusterName)
		}

		ginkgo.By("fetching an updated list of the configmaps from a cross-cluster context")
		configMaps, err = f.MultiClusterClientSet.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labels.Set{watchConfigMapNamespaceLabelKey: ns}).String(),
		})
		framework.ExpectNoError(err, "listing configmaps in all namespaces")

		ginkgo.By("deleting configmaps in single-cluster contexts")
		for _, configMap := range configMaps.Items {
			err := f.ClusterlessClientSet.Cluster(configMap.ClusterName).CoreV1().ConfigMaps(ns).Delete(context.TODO(), configMap.Name, metav1.DeleteOptions{
				Preconditions: &metav1.Preconditions{
					UID:             &configMap.UID,
					ResourceVersion: &configMap.ResourceVersion,
				},
			})
			framework.ExpectNoError(err, "deleting configmap %q in namespace %q for cluster %q", configMap.Name, configMap.Namespace, configMap.ClusterName)
		}
	})
})
