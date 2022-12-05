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

package v1beta1

import (
	"context"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	appsv1beta1client "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
)

// ControllerRevisionsClusterGetter has a method to return a ControllerRevisionClusterInterface.
// A group's cluster client should implement this interface.
type ControllerRevisionsClusterGetter interface {
	ControllerRevisions() ControllerRevisionClusterInterface
}

// ControllerRevisionClusterInterface can operate on ControllerRevisions across all clusters,
// or scope down to one cluster and return a ControllerRevisionsNamespacer.
type ControllerRevisionClusterInterface interface {
	Cluster(logicalcluster.Path) ControllerRevisionsNamespacer
	List(ctx context.Context, opts metav1.ListOptions) (*appsv1beta1.ControllerRevisionList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type controllerRevisionsClusterInterface struct {
	clientCache kcpclient.Cache[*appsv1beta1client.AppsV1beta1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *controllerRevisionsClusterInterface) Cluster(path logicalcluster.Path) ControllerRevisionsNamespacer {
	if path == logicalcluster.WildcardPath {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &controllerRevisionsNamespacer{clientCache: c.clientCache, path: path}
}

// List returns the entire collection of all ControllerRevisions across all clusters.
func (c *controllerRevisionsClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*appsv1beta1.ControllerRevisionList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).ControllerRevisions(metav1.NamespaceAll).List(ctx, opts)
}

// Watch begins to watch all ControllerRevisions across all clusters.
func (c *controllerRevisionsClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).ControllerRevisions(metav1.NamespaceAll).Watch(ctx, opts)
}

// ControllerRevisionsNamespacer can scope to objects within a namespace, returning a appsv1beta1client.ControllerRevisionInterface.
type ControllerRevisionsNamespacer interface {
	Namespace(string) appsv1beta1client.ControllerRevisionInterface
}

type controllerRevisionsNamespacer struct {
	clientCache kcpclient.Cache[*appsv1beta1client.AppsV1beta1Client]
	path        logicalcluster.Path
}

func (n *controllerRevisionsNamespacer) Namespace(namespace string) appsv1beta1client.ControllerRevisionInterface {
	return n.clientCache.ClusterOrDie(n.path).ControllerRevisions(namespace)
}
