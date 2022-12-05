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

	schedulingv1beta1 "k8s.io/api/scheduling/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	schedulingv1beta1client "k8s.io/client-go/kubernetes/typed/scheduling/v1beta1"
)

// PriorityClassesClusterGetter has a method to return a PriorityClassClusterInterface.
// A group's cluster client should implement this interface.
type PriorityClassesClusterGetter interface {
	PriorityClasses() PriorityClassClusterInterface
}

// PriorityClassClusterInterface can operate on PriorityClasses across all clusters,
// or scope down to one cluster and return a schedulingv1beta1client.PriorityClassInterface.
type PriorityClassClusterInterface interface {
	Cluster(logicalcluster.Path) schedulingv1beta1client.PriorityClassInterface
	List(ctx context.Context, opts metav1.ListOptions) (*schedulingv1beta1.PriorityClassList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type priorityClassesClusterInterface struct {
	clientCache kcpclient.Cache[*schedulingv1beta1client.SchedulingV1beta1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *priorityClassesClusterInterface) Cluster(path logicalcluster.Path) schedulingv1beta1client.PriorityClassInterface {
	if path == logicalcluster.WildcardPath {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return c.clientCache.ClusterOrDie(path).PriorityClasses()
}

// List returns the entire collection of all PriorityClasses across all clusters.
func (c *priorityClassesClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*schedulingv1beta1.PriorityClassList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).PriorityClasses().List(ctx, opts)
}

// Watch begins to watch all PriorityClasses across all clusters.
func (c *priorityClassesClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.WildcardPath).PriorityClasses().Watch(ctx, opts)
}
