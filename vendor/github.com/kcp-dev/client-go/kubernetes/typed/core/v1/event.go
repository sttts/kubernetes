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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

// EventsClusterGetter has a method to return a EventClusterInterface.
// A group's cluster client should implement this interface.
type EventsClusterGetter interface {
	Events() EventClusterInterface
}

// EventClusterInterface can operate on Events across all clusters,
// or scope down to one cluster and return a EventsNamespacer.
type EventClusterInterface interface {
	Cluster(logicalcluster.Path) EventsNamespacer
	List(ctx context.Context, opts metav1.ListOptions) (*corev1.EventList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type eventsClusterInterface struct {
	clientCache kcpclient.Cache[*corev1client.CoreV1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *eventsClusterInterface) Cluster(path logicalcluster.Path) EventsNamespacer {
	if path == logicalcluster.WildcardPath {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &eventsNamespacer{clientCache: c.clientCache, path: path}
}

// List returns the entire collection of all Events across all clusters.
func (c *eventsClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*corev1.EventList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).Events(metav1.NamespaceAll).List(ctx, opts)
}

// Watch begins to watch all Events across all clusters.
func (c *eventsClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).Events(metav1.NamespaceAll).Watch(ctx, opts)
}

// EventsNamespacer can scope to objects within a namespace, returning a corev1client.EventInterface.
type EventsNamespacer interface {
	Namespace(string) corev1client.EventInterface
}

type eventsNamespacer struct {
	clientCache kcpclient.Cache[*corev1client.CoreV1Client]
	path        logicalcluster.Path
}

func (n *eventsNamespacer) Namespace(namespace string) corev1client.EventInterface {
	return n.clientCache.ClusterOrDie(n.path).Events(namespace)
}
