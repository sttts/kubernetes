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

package apimachinery

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	cachetools "k8s.io/client-go/tools/cache"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/onsi/ginkgo"
)

const (
	watchConfigMapLabelKey          = "watch-this-configmap"
	watchConfigMapNamespaceLabelKey = "configmap-namespace"

	multipleWatchersLabelValueA   = "multiple-watchers-A"
	multipleWatchersLabelValueB   = "multiple-watchers-B"
	fromResourceVersionLabelValue = "from-resource-version"
	watchRestartedLabelValue      = "watch-closed-and-restarted"
	toBeChangedLabelValue         = "label-changed-and-restored"
)

var _ = SIGDescribe("KCP MultiCluster", func() {
	f := framework.NewDefaultFramework("watch")

	/*
		    Release: v1.11
		    Testname: watch-configmaps-with-multiple-watchers
		    Description: Ensure that multiple watchers are able to receive all add,
			update, and delete notifications on configmaps that match a label selector and do
			not receive notifications for configmaps which do not match that label selector.
	*/
	framework.ConformanceIt("should observe add, update, and delete watch notifications on configmaps", func() {
		if !framework.ProviderIs("kcp") {
			return // TODO: how to skip?
		}
		ns := f.Namespace.Name

		ginkgo.By("creating a watch on configmaps with label A")
		watchA, err := watchConfigMaps(f, "", multipleWatchersLabelValueA)
		framework.ExpectNoError(err, "failed to create a watch on configmaps with label: %s", multipleWatchersLabelValueA)

		ginkgo.By("creating a watch on configmaps with label B")
		watchB, err := watchConfigMaps(f, "", multipleWatchersLabelValueB)
		framework.ExpectNoError(err, "failed to create a watch on configmaps with label: %s", multipleWatchersLabelValueB)

		ginkgo.By("creating a watch on configmaps with label A or B")
		watchAB, err := watchConfigMaps(f, "", multipleWatchersLabelValueA, multipleWatchersLabelValueB)
		framework.ExpectNoError(err, "failed to create a watch on configmaps with label %s or %s", multipleWatchersLabelValueA, multipleWatchersLabelValueB)

		for i, c := range []clientset.Interface{f.ClientSet, f.SecondaryClientSet, f.TertiaryClientSet} {
			ginkgo.By(fmt.Sprintf("running with clientset %s", []string{"primary", "secondary", "tertiary"}[i]))
			if c != nil {
				testConfigMapA := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "e2e-watch-test-configmap-a",
						Labels: map[string]string{
							watchConfigMapNamespaceLabelKey: ns,
							watchConfigMapLabelKey:          multipleWatchersLabelValueA,
						},
					},
				}
				testConfigMapB := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "e2e-watch-test-configmap-b",
						Labels: map[string]string{
							watchConfigMapNamespaceLabelKey: ns,
							watchConfigMapLabelKey:          multipleWatchersLabelValueB,
						},
					},
				}

				ginkgo.By("creating a configmap with label A and ensuring the correct watchers observe the notification")
				testConfigMapA, err = c.CoreV1().ConfigMaps(ns).Create(context.TODO(), testConfigMapA, metav1.CreateOptions{})
				fmt.Printf("CREATE TYPEMETA: %#v\n", testConfigMapA.TypeMeta)
				framework.ExpectNoError(err, "failed to create a configmap with label %s in namespace: %s", multipleWatchersLabelValueA, ns)
				expectEvent(watchA, watch.Added, testConfigMapA)
				expectEvent(watchAB, watch.Added, testConfigMapA)
				expectNoEvent(watchB, watch.Added, testConfigMapA)

				ginkgo.By("modifying configmap A and ensuring the correct watchers observe the notification")
				testConfigMapA, err = updateConfigMap(c, ns, testConfigMapA.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "1")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace: %s", testConfigMapA.GetName(), ns)
				expectEvent(watchA, watch.Modified, testConfigMapA)
				expectEvent(watchAB, watch.Modified, testConfigMapA)
				expectNoEvent(watchB, watch.Modified, testConfigMapA)

				ginkgo.By("modifying configmap A again and ensuring the correct watchers observe the notification")
				testConfigMapA, err = updateConfigMap(c, ns, testConfigMapA.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "2")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace: %s", testConfigMapA.GetName(), ns)
				expectEvent(watchA, watch.Modified, testConfigMapA)
				expectEvent(watchAB, watch.Modified, testConfigMapA)
				expectNoEvent(watchB, watch.Modified, testConfigMapA)

				ginkgo.By("deleting configmap A and ensuring the correct watchers observe the notification")
				err = c.CoreV1().ConfigMaps(ns).Delete(context.TODO(), testConfigMapA.GetName(), metav1.DeleteOptions{})
				framework.ExpectNoError(err, "failed to delete configmap %s in namespace: %s", testConfigMapA.GetName(), ns)
				expectEvent(watchA, watch.Deleted, nil)
				expectEvent(watchAB, watch.Deleted, nil)
				expectNoEvent(watchB, watch.Deleted, nil)

				ginkgo.By("creating a configmap with label B and ensuring the correct watchers observe the notification")
				testConfigMapB, err = c.CoreV1().ConfigMaps(ns).Create(context.TODO(), testConfigMapB, metav1.CreateOptions{})
				framework.ExpectNoError(err, "failed to create configmap %s in namespace: %s", testConfigMapB, ns)
				expectEvent(watchB, watch.Added, testConfigMapB)
				expectEvent(watchAB, watch.Added, testConfigMapB)
				expectNoEvent(watchA, watch.Added, testConfigMapB)

				ginkgo.By("deleting configmap B and ensuring the correct watchers observe the notification")
				err = c.CoreV1().ConfigMaps(ns).Delete(context.TODO(), testConfigMapB.GetName(), metav1.DeleteOptions{})
				framework.ExpectNoError(err, "failed to delete configmap %s in namespace: %s", testConfigMapB.GetName(), ns)
				expectEvent(watchB, watch.Deleted, nil)
				expectEvent(watchAB, watch.Deleted, nil)
				expectNoEvent(watchA, watch.Deleted, nil)
			}
		}
	})

	/*
		    Release: v1.11
		    Testname: watch-configmaps-from-resource-version
		    Description: Ensure that a watch can be opened from a particular resource version
			in the past and only notifications happening after that resource version are observed.
	*/
	framework.ConformanceIt("should be able to start watching from a specific resource version", func() {
		ns := f.Namespace.Name

		for i, c := range []clientset.Interface{f.ClientSet, f.SecondaryClientSet, f.TertiaryClientSet} {
			ginkgo.By(fmt.Sprintf("running with clientset %s", []string{"primary", "secondary", "tertiary"}[i]))
			if c != nil {
				testConfigMap := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "e2e-watch-test-resource-version",
						Labels: map[string]string{
							watchConfigMapNamespaceLabelKey: ns,
							watchConfigMapLabelKey:          fromResourceVersionLabelValue,
						},
					},
				}

				ginkgo.By("creating a new configmap")
				testConfigMap, err := c.CoreV1().ConfigMaps(ns).Create(context.TODO(), testConfigMap, metav1.CreateOptions{})
				framework.ExpectNoError(err, "failed to create configmap %s in namespace: %s", testConfigMap.GetName(), ns)

				ginkgo.By("modifying the configmap once")
				_, err = updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "1")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace: %s", testConfigMap.GetName(), ns)
				listAfterFirstUpdate, err := f.MultiClusterClientSet.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
				framework.ExpectNoError(err, "failed to list configmaps while getting current resource version")
				fmt.Printf("GOT RESOURCE VERSION FROM LIST: %S", listAfterFirstUpdate.ListMeta.ResourceVersion)

				ginkgo.By("modifying the configmap a second time")
				testConfigMapSecondUpdate, err := updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "2")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace %s a second time", testConfigMap.GetName(), ns)

				ginkgo.By("deleting the configmap")
				err = c.CoreV1().ConfigMaps(ns).Delete(context.TODO(), testConfigMap.GetName(), metav1.DeleteOptions{})
				framework.ExpectNoError(err, "failed to delete configmap %s in namespace: %s", testConfigMap.GetName(), ns)

				ginkgo.By("creating a watch on configmaps from the resource version returned by the first update")
				testWatch, err := watchConfigMaps(f, listAfterFirstUpdate.ListMeta.ResourceVersion, fromResourceVersionLabelValue)
				framework.ExpectNoError(err, "failed to create a watch on configmaps from the resource version %s returned by a list after the first update", listAfterFirstUpdate.ListMeta.ResourceVersion)

				ginkgo.By("Expecting to observe notifications for all changes to the configmap after the first update")
				expectEvent(testWatch, watch.Modified, testConfigMapSecondUpdate)
				expectEvent(testWatch, watch.Deleted, nil)
			}
		}
	})

	/*
		    Release: v1.11
		    Testname: watch-configmaps-closed-and-restarted
		    Description: Ensure that a watch can be reopened from the last resource version
			observed by the previous watch, and it will continue delivering notifications from
			that point in time.
	*/
	framework.ConformanceIt("should be able to restart watching from the last resource version observed by the previous watch", func() {
		ns := f.Namespace.Name

		for i, c := range []clientset.Interface{f.ClientSet, f.SecondaryClientSet, f.TertiaryClientSet} {
			ginkgo.By(fmt.Sprintf("running with clientset %s", []string{"primary", "secondary", "tertiary"}[i]))
			if c != nil {
				configMapName := "e2e-watch-test-watch-closed"
				testConfigMap := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: configMapName,
						Labels: map[string]string{
							watchConfigMapNamespaceLabelKey: ns,
							watchConfigMapLabelKey:          watchRestartedLabelValue,
						},
					},
				}

				ginkgo.By("creating a watch on configmaps")
				testWatchBroken, err := watchConfigMaps(f, "", watchRestartedLabelValue)
				framework.ExpectNoError(err, "failed to create a watch on configmap with label: %s", watchRestartedLabelValue)

				ginkgo.By("creating a new configmap")
				testConfigMap, err = c.CoreV1().ConfigMaps(ns).Create(context.TODO(), testConfigMap, metav1.CreateOptions{})
				framework.ExpectNoError(err, "failed to create configmap %s in namespace: %s", configMapName, ns)

				ginkgo.By("modifying the configmap once")
				_, err = updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "1")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace: %s", configMapName, ns)

				ginkgo.By("closing the watch once it receives two notifications")
				expectEvent(testWatchBroken, watch.Added, testConfigMap)
				lastEvent, ok := waitForEvent(testWatchBroken, watch.Modified, nil, 1*time.Minute)
				if !ok {
					framework.Failf("Timed out waiting for second watch notification")
				}
				testWatchBroken.Stop()

				ginkgo.By("modifying the configmap a second time, while the watch is closed")
				testConfigMapSecondUpdate, err := updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "2")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace %s a second time", configMapName, ns)

				ginkgo.By("creating a new watch on configmaps from the last resource version observed by the first watch")
				lastEventConfigMap, ok := lastEvent.Object.(*v1.ConfigMap)
				if !ok {
					framework.Failf("Expected last notification to refer to a configmap but got: %v", lastEvent)
				}
				testWatchRestarted, err := watchConfigMaps(f, lastEventConfigMap.ObjectMeta.ResourceVersion, watchRestartedLabelValue)
				framework.ExpectNoError(err, "failed to create a new watch on configmaps from the last resource version %s observed by the first watch", lastEventConfigMap.ObjectMeta.ResourceVersion)

				ginkgo.By("deleting the configmap")
				err = c.CoreV1().ConfigMaps(ns).Delete(context.TODO(), testConfigMap.GetName(), metav1.DeleteOptions{})
				framework.ExpectNoError(err, "failed to delete configmap %s in namespace: %s", configMapName, ns)

				ginkgo.By("Expecting to observe notifications for all changes to the configmap since the first watch closed")
				expectEvent(testWatchRestarted, watch.Modified, testConfigMapSecondUpdate)
				expectEvent(testWatchRestarted, watch.Deleted, nil)
			}
		}
	})

	/*
		    Release: v1.11
		    Testname: watch-configmaps-label-changed
		    Description: Ensure that a watched object stops meeting the requirements of
			a watch's selector, the watch will observe a delete, and will not observe
			notifications for that object until it meets the selector's requirements again.
	*/
	framework.ConformanceIt("should observe an object deletion if it stops meeting the requirements of the selector", func() {
		ns := f.Namespace.Name

		ginkgo.By("creating a watch on configmaps with a certain label")
		testWatch, err := watchConfigMaps(f, "", toBeChangedLabelValue)
		framework.ExpectNoError(err, "failed to create a watch on configmap with label: %s", toBeChangedLabelValue)

		for i, c := range []clientset.Interface{f.ClientSet, f.SecondaryClientSet, f.TertiaryClientSet} {
			ginkgo.By(fmt.Sprintf("running with clientset %s", []string{"primary", "secondary", "tertiary"}[i]))
			if c != nil {
				configMapName := "e2e-watch-test-label-changed"
				testConfigMap := &v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: configMapName,
						Labels: map[string]string{
							watchConfigMapNamespaceLabelKey: ns,
							watchConfigMapLabelKey:          toBeChangedLabelValue,
						},
					},
				}

				ginkgo.By("creating a new configmap")
				testConfigMap, err = c.CoreV1().ConfigMaps(ns).Create(context.TODO(), testConfigMap, metav1.CreateOptions{})
				framework.ExpectNoError(err, "failed to create configmap %s in namespace: %s", configMapName, ns)

				ginkgo.By("modifying the configmap once")
				testConfigMapFirstUpdate, err := updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "1")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace: %s", configMapName, ns)

				ginkgo.By("changing the label value of the configmap")
				_, err = updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					cm.ObjectMeta.Labels[watchConfigMapLabelKey] = "wrong-value"
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace %s by changing label value", configMapName, ns)

				ginkgo.By("Expecting to observe a delete notification for the watched object")
				expectEvent(testWatch, watch.Added, testConfigMap)
				expectEvent(testWatch, watch.Modified, testConfigMapFirstUpdate)
				expectEvent(testWatch, watch.Deleted, nil)

				ginkgo.By("modifying the configmap a second time")
				testConfigMapSecondUpdate, err := updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "2")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace %s a second time", configMapName, ns)

				ginkgo.By("Expecting not to observe a notification because the object no longer meets the selector's requirements")
				expectNoEvent(testWatch, watch.Modified, testConfigMapSecondUpdate)

				ginkgo.By("changing the label value of the configmap back")
				testConfigMapLabelRestored, err := updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					cm.ObjectMeta.Labels[watchConfigMapLabelKey] = toBeChangedLabelValue
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace %s by changing label value back", configMapName, ns)

				ginkgo.By("modifying the configmap a third time")
				testConfigMapThirdUpdate, err := updateConfigMap(c, ns, testConfigMap.GetName(), func(cm *v1.ConfigMap) {
					setConfigMapData(cm, "mutation", "3")
				})
				framework.ExpectNoError(err, "failed to update configmap %s in namespace %s a third time", configMapName, ns)

				ginkgo.By("deleting the configmap")
				err = c.CoreV1().ConfigMaps(ns).Delete(context.TODO(), testConfigMap.GetName(), metav1.DeleteOptions{})
				framework.ExpectNoError(err, "failed to delete configmap %s in namespace: %s", configMapName, ns)

				ginkgo.By("Expecting to observe an add notification for the watched object when the label value was restored")
				expectEvent(testWatch, watch.Added, testConfigMapLabelRestored)
				expectEvent(testWatch, watch.Modified, testConfigMapThirdUpdate)
				expectEvent(testWatch, watch.Deleted, nil)
			}
		}
	})

	/*
	   Release: v1.15
	   Testname: watch-consistency
	   Description: Ensure that concurrent watches are consistent with each other by initiating an additional watch
	   for events received from the first watch, initiated at the resource version of the event, and checking that all
	   resource versions of all events match. Events are produced from writes on a background goroutine.
	*/
	framework.ConformanceIt("should receive events on concurrent watches in same order", func() {
		ns := f.Namespace.Name

		iterations := 35
		totalIterations := 0

		ginkgo.By("getting a starting resourceVersion")
		configmaps, err := f.MultiClusterClientSet.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
		framework.ExpectNoError(err, "Failed to list configmaps")
		resourceVersion := configmaps.ResourceVersion

		ginkgo.By("starting a background goroutine to produce watch events")
		done := sync.WaitGroup{}
		stopc := make(chan struct{})
		for i, c := range []clientset.Interface{f.ClientSet, f.SecondaryClientSet, f.TertiaryClientSet} {
			ginkgo.By(fmt.Sprintf("running with clientset %s", []string{"primary", "secondary", "tertiary"}[i]))
			if c != nil {
				done.Add(1)
				totalIterations += iterations
				go func(c clientset.Interface) {
					defer ginkgo.GinkgoRecover()
					defer done.Done()
					produceConfigMapEvents(c, ns, iterations, stopc, 5*time.Millisecond)
				}(c)
			}
		}

		listWatcher := &cachetools.ListWatch{
			WatchFunc: func(listOptions metav1.ListOptions) (watch.Interface, error) {
				opts := listOptions.DeepCopy()
				opts.LabelSelector = metav1.FormatLabelSelector(&metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						// we can't do namespaced list/watch cross-cluster, and need this to be unique for parallelism...
						{
							Key:      watchConfigMapNamespaceLabelKey,
							Operator: metav1.LabelSelectorOpIn,
							Values:   []string{f.Namespace.Name},
						},
					},
				})
				return f.MultiClusterClientSet.CoreV1().ConfigMaps("").Watch(context.TODO(), listOptions)
			},
		}

		done.Wait()

		ginkgo.By("creating watches starting from each resource version of the events produced and verifying they all receive resource versions from a given cluster in the same order")
		simpleRvFor := func(e *v1.ConfigMap) int64 {
			complexRv := ShardedResourceVersions{}
			framework.ExpectNoError(complexRv.Decode(e.ResourceVersion), "decoding complex RV from watch event")
			simpleRv := int64(-1)
			for _, item := range complexRv.ResourceVersions {
				if item.Identifier == e.ObjectMeta.ClusterName {
					simpleRv = item.ResourceVersion
					break
				}
			}
			if simpleRv == -1 {
				framework.Failf("could not find resource version for cluster %q in complex RV %#v", e.ObjectMeta.ClusterName, complexRv)
			}
			return simpleRv
		}
		var resourceVersions []map[string][]int64
		wcs := []watch.Interface{}
		for i := 0; i < totalIterations; i++ {
			wc, err := watchtools.NewRetryWatcher(resourceVersion, listWatcher)
			framework.ExpectNoError(err, "Failed to watch configmaps in the namespace %s", ns)
			wcs = append(wcs, wc)
			e := waitForNextConfigMapEvent(wcs[0])
			resourceVersions = append(resourceVersions, map[string][]int64{
				e.ObjectMeta.ClusterName: {simpleRvFor(e)},
			})
			resourceVersion = e.ResourceVersion
			for j, wc := range wcs[1:] {
				e := waitForNextConfigMapEvent(wc)
				resourceVersions[j][e.ObjectMeta.ClusterName] = append(resourceVersions[j][e.ObjectMeta.ClusterName], simpleRvFor(e))
			}
		}
		for i, allRvs := range resourceVersions[1:] {
			for cluster, rvs := range allRvs {
				reference, seen := resourceVersions[0][cluster]
				if !seen {
					framework.Failf("watcher %d saw events for cluster %q, but the reference watcher did not!", i+1, cluster)
				}
				startingIndex := -1
				for j, item := range reference {
					if item == rvs[0] {
						startingIndex = j
					}
				}
				if startingIndex == -1 {
					fmt.Println(rvs)
					fmt.Println(reference)
					framework.Failf("watcher %d started seeing events for cluster %q at rv %q, but the reference watcher never saw that version!", i+1, cluster, rvs[0])
				}
				if diff := cmp.Diff(rvs, reference[startingIndex:startingIndex+len(rvs)]); diff != "" {
					framework.Failf("watcher %d got incorrect ordering for rvs for cluster %q: %v", i+1, cluster, diff)
				}

				if !sort.SliceIsSorted(rvs, func(i, j int) bool {
					return rvs[i] < rvs[j]
				}) {
					framework.Failf("watcher %d got out-of-order rvs for cluster %q: %v", i+1, cluster, rvs)
				}
			}
		}
		close(stopc)
		for _, wc := range wcs {
			wc.Stop()
		}
	})
})

func watchConfigMaps(f *framework.Framework, resourceVersion string, labels ...string) (watch.Interface, error) {
	c := f.MultiClusterClientSet
	opts := metav1.ListOptions{
		ResourceVersion: resourceVersion,
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      watchConfigMapLabelKey,
					Operator: metav1.LabelSelectorOpIn,
					Values:   labels,
				},
				// we can't do namespaced list/watch cross-cluster, and need this to be unique for parallelism...
				{
					Key:      watchConfigMapNamespaceLabelKey,
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{f.Namespace.Name},
				},
			},
		}),
	}
	return c.CoreV1().ConfigMaps("").Watch(context.TODO(), opts)
}

func int64ptr(i int) *int64 {
	i64 := int64(i)
	return &i64
}

func setConfigMapData(cm *v1.ConfigMap, key, value string) {
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}
	cm.Data[key] = value
}

func expectEvent(w watch.Interface, eventType watch.EventType, object runtime.Object) {
	if event, ok := waitForEvent(w, eventType, object, 1*time.Minute); !ok {
		framework.Failf("Timed out waiting for expected watch notification: %v", event)
	}
}

func expectNoEvent(w watch.Interface, eventType watch.EventType, object runtime.Object) {
	if event, ok := waitForEvent(w, eventType, object, 10*time.Second); ok {
		framework.Failf("Unexpected watch notification observed: %v", event)
	}
}

func waitForEvent(w watch.Interface, expectType watch.EventType, expectObject runtime.Object, duration time.Duration) (watch.Event, bool) {
	stopTimer := time.NewTimer(duration)
	defer stopTimer.Stop()
	if expectObject != nil {
		// simple (single-cluster) RV can't be compared to those on watch events without racing, so we remove it ...
		expectObject = expectObject.DeepCopyObject()
		if obj, ok := expectObject.(metav1.Common); ok {
			obj.SetResourceVersion("")
		} else {
			framework.Failf("was passed %T and could not reset RV", expectObject)
		}
	}
	for {
		select {
		case actual, ok := <-w.ResultChan():
			copied := actual.DeepCopy()
			if ok {
				framework.Logf("Got : %v %v", actual.Type, actual.Object)
				if obj, ok := actual.Object.(metav1.Common); ok {
					obj.SetResourceVersion("")
				} else if actual.Object != nil {
					framework.Failf("got %T and could not reset RV", actual.Object)
				}
			} else {
				framework.Failf("Watch closed unexpectedly")
			}
			if expectType == actual.Type && (expectObject == nil || apiequality.Semantic.DeepEqual(expectObject, actual.Object)) {
				return *copied, true
			} else if expectType != actual.Type {
				fmt.Printf("INCORRECT TYPE: %v != %v\n", expectType, actual.Type)
			} else if expectObject == nil && actual.Object != nil {
				fmt.Printf("EXPECTED NIL, GOT: %#v\n", actual.Object)
			} else {
				fmt.Printf("DIFF: %s\n", cmp.Diff(expectObject, actual.Object))
			}
		case <-stopTimer.C:
			expected := watch.Event{
				Type:   expectType,
				Object: expectObject,
			}
			return expected, false
		}
	}
}

func waitForNextConfigMapEvent(watch watch.Interface) *v1.ConfigMap {
	select {
	case event, ok := <-watch.ResultChan():
		if !ok {
			framework.Failf("Watch closed unexpectedly")
		}
		if configMap, ok := event.Object.(*v1.ConfigMap); ok {
			return configMap
		}
		framework.Failf("expected config map, got %T", event.Object)
	case <-time.After(10 * time.Second):
		framework.Failf("timed out waiting for watch event")
	}
	return nil // should never happen
}

const (
	createEvent = iota
	updateEvent
	deleteEvent
)

func produceConfigMapEvents(c clientset.Interface, ns string, iterations int, stopc <-chan struct{}, minWaitBetweenEvents time.Duration) {
	name := func(i int) string {
		return fmt.Sprintf("cm-%d", i)
	}

	existing := []int{}
	tc := time.NewTicker(minWaitBetweenEvents)
	defer tc.Stop()
	var creates, updates, deletes int
	i := 0
	for x := 0; x <iterations; x++ {
		<-tc.C
		op := rand.Intn(3)
		if len(existing) == 0 {
			op = createEvent
		}

		cm := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					watchConfigMapNamespaceLabelKey: ns,
				},
			},
		}
		switch op {
		case createEvent:
			toCreate := cm.DeepCopy()
			toCreate.Name = name(i)
			_, err := c.CoreV1().ConfigMaps(ns).Create(context.TODO(), toCreate, metav1.CreateOptions{})
			framework.ExpectNoError(err, "Failed to create configmap %s in namespace %s", toCreate.Name, ns)
			existing = append(existing, i)
			i++
			creates++
		case updateEvent:
			idx := rand.Intn(len(existing))
			toUpdate := cm.DeepCopy()
			toUpdate.Name = name(existing[idx])
			toUpdate.ObjectMeta.Labels["mutated"] = strconv.Itoa(updates)
			_, err := c.CoreV1().ConfigMaps(ns).Update(context.TODO(), toUpdate, metav1.UpdateOptions{})
			framework.ExpectNoError(err, "Failed to update configmap %s in namespace %s", toUpdate.Name, ns)
			updates++
		case deleteEvent:
			idx := rand.Intn(len(existing))
			err := c.CoreV1().ConfigMaps(ns).Delete(context.TODO(), name(existing[idx]), metav1.DeleteOptions{})
			framework.ExpectNoError(err, "Failed to delete configmap %s in namespace %s", name(existing[idx]), ns)
			existing = append(existing[:idx], existing[idx+1:]...)
			deletes++
		default:
			framework.Failf("Unsupported event operation: %d", op)
		}
		select {
		case <-stopc:
			return
		default:
		}
	}
	fmt.Printf("FINISHED EVENTS %d/%d/%d\n", creates, updates, deletes)
}

type updateConfigMapFn func(cm *v1.ConfigMap)

func updateConfigMap(c clientset.Interface, ns, name string, update updateConfigMapFn) (*v1.ConfigMap, error) {
	var cm *v1.ConfigMap
	pollErr := wait.PollImmediate(2*time.Second, 1*time.Minute, func() (bool, error) {
		var err error
		if cm, err = c.CoreV1().ConfigMaps(ns).Get(context.TODO(), name, metav1.GetOptions{}); err != nil {
			return false, err
		}
		update(cm)
		if cm, err = c.CoreV1().ConfigMaps(ns).Update(context.TODO(), cm, metav1.UpdateOptions{}); err == nil {
			return true, nil
		}
		// Only retry update on conflict
		if !apierrors.IsConflict(err) {
			return false, err
		}
		return false, nil
	})
	return cm, pollErr
}


// TODO: move this and the tests into kcp...

// ShardedResourceVersions are what a client passes to the resourceVersion
// query parameter to initiate a LIST or WATCH across shards at a particular
// point in time.
type ShardedResourceVersions struct {
	// ShardResourceVersion is the version at which the list of shards
	// that needed to be queried was resolved.
	ShardResourceVersion int64 `json:"srv"`
	// ResourceVersions hold versions for individual shards being queried.
	ResourceVersions []ShardedResourceVersion `json:"rvs"`
}

type ShardedResourceVersion struct {
	// Identifier is used to locate credentials for this shard
	Identifier      string `json:"id"`
	ResourceVersion int64  `json:"rv"`
}

// Decode decodes a resource version into state
func (s *ShardedResourceVersions) Decode(encoded string) error {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("invalid sharded resource version encoding: %w", err)
	}
	if err := json.Unmarshal(raw, s); err != nil {
		return fmt.Errorf("invalid sharded resource version serialization: %w", err)
	}
	return nil
}