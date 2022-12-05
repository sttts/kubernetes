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

	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	extensionsv1beta1client "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

// DaemonSetsClusterGetter has a method to return a DaemonSetClusterInterface.
// A group's cluster client should implement this interface.
type DaemonSetsClusterGetter interface {
	DaemonSets() DaemonSetClusterInterface
}

// DaemonSetClusterInterface can operate on DaemonSets across all clusters,
// or scope down to one cluster and return a DaemonSetsNamespacer.
type DaemonSetClusterInterface interface {
	Cluster(logicalcluster.Path) DaemonSetsNamespacer
	List(ctx context.Context, opts metav1.ListOptions) (*extensionsv1beta1.DaemonSetList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type daemonSetsClusterInterface struct {
	clientCache kcpclient.Cache[*extensionsv1beta1client.ExtensionsV1beta1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *daemonSetsClusterInterface) Cluster(path logicalcluster.Path) DaemonSetsNamespacer {
	if path == logicalcluster.WildcardPath {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &daemonSetsNamespacer{clientCache: c.clientCache, path: path}
}

// List returns the entire collection of all DaemonSets across all clusters.
func (c *daemonSetsClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*extensionsv1beta1.DaemonSetList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).DaemonSets(metav1.NamespaceAll).List(ctx, opts)
}

// Watch begins to watch all DaemonSets across all clusters.
func (c *daemonSetsClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).DaemonSets(metav1.NamespaceAll).Watch(ctx, opts)
}

// DaemonSetsNamespacer can scope to objects within a namespace, returning a extensionsv1beta1client.DaemonSetInterface.
type DaemonSetsNamespacer interface {
	Namespace(string) extensionsv1beta1client.DaemonSetInterface
}

type daemonSetsNamespacer struct {
	clientCache kcpclient.Cache[*extensionsv1beta1client.ExtensionsV1beta1Client]
	path        logicalcluster.Path
}

func (n *daemonSetsNamespacer) Namespace(namespace string) extensionsv1beta1client.DaemonSetInterface {
	return n.clientCache.ClusterOrDie(n.path).DaemonSets(namespace)
}
