//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The KCP Authors.

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

// Code generated by kcp code-generator. DO NOT EDIT.

package v1beta2

import (
	"context"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	appsv1beta2client "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
)

// StatefulSetsClusterGetter has a method to return a StatefulSetClusterInterface.
// A group's cluster client should implement this interface.
type StatefulSetsClusterGetter interface {
	StatefulSets() StatefulSetClusterInterface
}

// StatefulSetClusterInterface can operate on StatefulSets across all clusters,
// or scope down to one cluster and return a StatefulSetsNamespacer.
type StatefulSetClusterInterface interface {
	Cluster(logicalcluster.Path) StatefulSetsNamespacer
	List(ctx context.Context, opts metav1.ListOptions) (*appsv1beta2.StatefulSetList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type statefulSetsClusterInterface struct {
	clientCache kcpclient.Cache[*appsv1beta2client.AppsV1beta2Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *statefulSetsClusterInterface) Cluster(name logicalcluster.Path) StatefulSetsNamespacer {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &statefulSetsNamespacer{clientCache: c.clientCache, name: name}
}

// List returns the entire collection of all StatefulSets across all clusters.
func (c *statefulSetsClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*appsv1beta2.StatefulSetList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).StatefulSets(metav1.NamespaceAll).List(ctx, opts)
}

// Watch begins to watch all StatefulSets across all clusters.
func (c *statefulSetsClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).StatefulSets(metav1.NamespaceAll).Watch(ctx, opts)
}

// StatefulSetsNamespacer can scope to objects within a namespace, returning a appsv1beta2client.StatefulSetInterface.
type StatefulSetsNamespacer interface {
	Namespace(string) appsv1beta2client.StatefulSetInterface
}

type statefulSetsNamespacer struct {
	clientCache kcpclient.Cache[*appsv1beta2client.AppsV1beta2Client]
	name        logicalcluster.Path
}

func (n *statefulSetsNamespacer) Namespace(namespace string) appsv1beta2client.StatefulSetInterface {
	return n.clientCache.ClusterOrDie(n.name).StatefulSets(namespace)
}
