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
	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"
	"github.com/kcp-dev/logicalcluster/v3"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	appsv1beta2listers "k8s.io/client-go/listers/apps/v1beta2"
	"k8s.io/client-go/tools/cache"
)

// DeploymentClusterLister can list Deployments across all workspaces, or scope down to a DeploymentLister for one workspace.
// All objects returned here must be treated as read-only.
type DeploymentClusterLister interface {
	// List lists all Deployments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*appsv1beta2.Deployment, err error)
	// Cluster returns a lister that can list and get Deployments in one workspace.
	Cluster(cluster logicalcluster.Name) appsv1beta2listers.DeploymentLister
	DeploymentClusterListerExpansion
}

type deploymentClusterLister struct {
	indexer cache.Indexer
}

// NewDeploymentClusterLister returns a new DeploymentClusterLister.
// We assume that the indexer:
// - is fed by a cross-workspace LIST+WATCH
// - uses kcpcache.MetaClusterNamespaceKeyFunc as the key function
// - has the kcpcache.ClusterIndex as an index
// - has the kcpcache.ClusterAndNamespaceIndex as an index
func NewDeploymentClusterLister(indexer cache.Indexer) *deploymentClusterLister {
	return &deploymentClusterLister{indexer: indexer}
}

// List lists all Deployments in the indexer across all workspaces.
func (s *deploymentClusterLister) List(selector labels.Selector) (ret []*appsv1beta2.Deployment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*appsv1beta2.Deployment))
	})
	return ret, err
}

// Cluster scopes the lister to one workspace, allowing users to list and get Deployments.
func (s *deploymentClusterLister) Cluster(cluster logicalcluster.Name) appsv1beta2listers.DeploymentLister {
	return &deploymentLister{indexer: s.indexer, cluster: cluster}
}

// deploymentLister implements the appsv1beta2listers.DeploymentLister interface.
type deploymentLister struct {
	indexer cache.Indexer
	cluster logicalcluster.Name
}

// List lists all Deployments in the indexer for a workspace.
func (s *deploymentLister) List(selector labels.Selector) (ret []*appsv1beta2.Deployment, err error) {
	err = kcpcache.ListAllByCluster(s.indexer, s.cluster, selector, func(i interface{}) {
		ret = append(ret, i.(*appsv1beta2.Deployment))
	})
	return ret, err
}

// Deployments returns an object that can list and get Deployments in one namespace.
func (s *deploymentLister) Deployments(namespace string) appsv1beta2listers.DeploymentNamespaceLister {
	return &deploymentNamespaceLister{indexer: s.indexer, cluster: s.cluster, namespace: namespace}
}

// deploymentNamespaceLister implements the appsv1beta2listers.DeploymentNamespaceLister interface.
type deploymentNamespaceLister struct {
	indexer   cache.Indexer
	cluster   logicalcluster.Name
	namespace string
}

// List lists all Deployments in the indexer for a given workspace and namespace.
func (s *deploymentNamespaceLister) List(selector labels.Selector) (ret []*appsv1beta2.Deployment, err error) {
	err = kcpcache.ListAllByClusterAndNamespace(s.indexer, s.cluster, s.namespace, selector, func(i interface{}) {
		ret = append(ret, i.(*appsv1beta2.Deployment))
	})
	return ret, err
}

// Get retrieves the Deployment from the indexer for a given workspace, namespace and name.
func (s *deploymentNamespaceLister) Get(name string) (*appsv1beta2.Deployment, error) {
	key := kcpcache.ToClusterAwareKey(s.cluster.String(), s.namespace, name)
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(appsv1beta2.Resource("Deployment"), name)
	}
	return obj.(*appsv1beta2.Deployment), nil
}
