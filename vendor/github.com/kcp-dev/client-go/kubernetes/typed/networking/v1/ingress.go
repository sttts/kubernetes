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

package v1

import (
	"context"

	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	networkingv1client "k8s.io/client-go/kubernetes/typed/networking/v1"
)

// IngressesClusterGetter has a method to return a IngressClusterInterface.
// A group's cluster client should implement this interface.
type IngressesClusterGetter interface {
	Ingresses() IngressClusterInterface
}

// IngressClusterInterface can operate on Ingresses across all clusters,
// or scope down to one cluster and return a IngressesNamespacer.
type IngressClusterInterface interface {
	Cluster(logicalcluster.Path) IngressesNamespacer
	List(ctx context.Context, opts metav1.ListOptions) (*networkingv1.IngressList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type ingressesClusterInterface struct {
	clientCache kcpclient.Cache[*networkingv1client.NetworkingV1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *ingressesClusterInterface) Cluster(name logicalcluster.Path) IngressesNamespacer {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &ingressesNamespacer{clientCache: c.clientCache, name: name}
}

// List returns the entire collection of all Ingresses across all clusters.
func (c *ingressesClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*networkingv1.IngressList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).Ingresses(metav1.NamespaceAll).List(ctx, opts)
}

// Watch begins to watch all Ingresses across all clusters.
func (c *ingressesClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).Ingresses(metav1.NamespaceAll).Watch(ctx, opts)
}

// IngressesNamespacer can scope to objects within a namespace, returning a networkingv1client.IngressInterface.
type IngressesNamespacer interface {
	Namespace(string) networkingv1client.IngressInterface
}

type ingressesNamespacer struct {
	clientCache kcpclient.Cache[*networkingv1client.NetworkingV1Client]
	name        logicalcluster.Path
}

func (n *ingressesNamespacer) Namespace(namespace string) networkingv1client.IngressInterface {
	return n.clientCache.ClusterOrDie(n.name).Ingresses(namespace)
}
