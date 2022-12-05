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

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	storagev1client "k8s.io/client-go/kubernetes/typed/storage/v1"
)

// CSIStorageCapacitiesClusterGetter has a method to return a CSIStorageCapacityClusterInterface.
// A group's cluster client should implement this interface.
type CSIStorageCapacitiesClusterGetter interface {
	CSIStorageCapacities() CSIStorageCapacityClusterInterface
}

// CSIStorageCapacityClusterInterface can operate on CSIStorageCapacities across all clusters,
// or scope down to one cluster and return a CSIStorageCapacitiesNamespacer.
type CSIStorageCapacityClusterInterface interface {
	Cluster(logicalcluster.Path) CSIStorageCapacitiesNamespacer
	List(ctx context.Context, opts metav1.ListOptions) (*storagev1.CSIStorageCapacityList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type cSIStorageCapacitiesClusterInterface struct {
	clientCache kcpclient.Cache[*storagev1client.StorageV1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *cSIStorageCapacitiesClusterInterface) Cluster(name logicalcluster.Path) CSIStorageCapacitiesNamespacer {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &cSIStorageCapacitiesNamespacer{clientCache: c.clientCache, name: name}
}

// List returns the entire collection of all CSIStorageCapacities across all clusters.
func (c *cSIStorageCapacitiesClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*storagev1.CSIStorageCapacityList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).CSIStorageCapacities(metav1.NamespaceAll).List(ctx, opts)
}

// Watch begins to watch all CSIStorageCapacities across all clusters.
func (c *cSIStorageCapacitiesClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).CSIStorageCapacities(metav1.NamespaceAll).Watch(ctx, opts)
}

// CSIStorageCapacitiesNamespacer can scope to objects within a namespace, returning a storagev1client.CSIStorageCapacityInterface.
type CSIStorageCapacitiesNamespacer interface {
	Namespace(string) storagev1client.CSIStorageCapacityInterface
}

type cSIStorageCapacitiesNamespacer struct {
	clientCache kcpclient.Cache[*storagev1client.StorageV1Client]
	name        logicalcluster.Path
}

func (n *cSIStorageCapacitiesNamespacer) Namespace(namespace string) storagev1client.CSIStorageCapacityInterface {
	return n.clientCache.ClusterOrDie(n.name).CSIStorageCapacities(namespace)
}
