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
	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"
	"github.com/kcp-dev/logicalcluster/v2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// ReplicationControllerClusterLister can list ReplicationControllers across all workspaces, or scope down to a ReplicationControllerLister for one workspace.
type ReplicationControllerClusterLister interface {
	List(selector labels.Selector) (ret []*corev1.ReplicationController, err error)
	Cluster(cluster logicalcluster.Name) corev1listers.ReplicationControllerLister
}

type replicationControllerClusterLister struct {
	indexer cache.Indexer
}

// NewReplicationControllerClusterLister returns a new ReplicationControllerClusterLister.
func NewReplicationControllerClusterLister(indexer cache.Indexer) *replicationControllerClusterLister {
	return &replicationControllerClusterLister{indexer: indexer}
}

// List lists all ReplicationControllers in the indexer across all workspaces.
func (s *replicationControllerClusterLister) List(selector labels.Selector) (ret []*corev1.ReplicationController, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*corev1.ReplicationController))
	})
	return ret, err
}

// Cluster scopes the lister to one workspace, allowing users to list and get ReplicationControllers.
func (s *replicationControllerClusterLister) Cluster(cluster logicalcluster.Name) corev1listers.ReplicationControllerLister {
	return &replicationControllerLister{indexer: s.indexer, cluster: cluster}
}

// replicationControllerLister implements the corev1listers.ReplicationControllerLister interface.
type replicationControllerLister struct {
	indexer cache.Indexer
	cluster logicalcluster.Name
}

// List lists all ReplicationControllers in the indexer for a workspace.
func (s *replicationControllerLister) List(selector labels.Selector) (ret []*corev1.ReplicationController, err error) {
	selectAll := selector == nil || selector.Empty()

	list, err := s.indexer.ByIndex(kcpcache.ClusterIndexName, kcpcache.ClusterIndexKey(s.cluster))
	if err != nil {
		return nil, err
	}

	for i := range list {
		obj := list[i].(*corev1.ReplicationController)
		if selectAll {
			ret = append(ret, obj)
		} else {
			if selector.Matches(labels.Set(obj.GetLabels())) {
				ret = append(ret, obj)
			}
		}
	}

	return ret, err
}

// ReplicationControllers returns an object that can list and get ReplicationControllers in one namespace.
func (s *replicationControllerLister) ReplicationControllers(namespace string) corev1listers.ReplicationControllerNamespaceLister {
	return &replicationControllerNamespaceLister{indexer: s.indexer, cluster: s.cluster, namespace: namespace}
}

// replicationControllerNamespaceLister implements the corev1listers.ReplicationControllerNamespaceLister interface.
type replicationControllerNamespaceLister struct {
	indexer   cache.Indexer
	cluster   logicalcluster.Name
	namespace string
}

// List lists all ReplicationControllers in the indexer for a given workspace and namespace.
func (s *replicationControllerNamespaceLister) List(selector labels.Selector) (ret []*corev1.ReplicationController, err error) {
	selectAll := selector == nil || selector.Empty()

	var list []interface{}
	if s.namespace == metav1.NamespaceAll {
		list, err = s.indexer.ByIndex(kcpcache.ClusterIndexName, kcpcache.ClusterIndexKey(s.cluster))
	} else {
		list, err = s.indexer.ByIndex(kcpcache.ClusterAndNamespaceIndexName, kcpcache.ClusterAndNamespaceIndexKey(s.cluster, s.namespace))
	}
	if err != nil {
		return nil, err
	}

	for i := range list {
		obj := list[i].(*corev1.ReplicationController)
		if selectAll {
			ret = append(ret, obj)
		} else {
			if selector.Matches(labels.Set(obj.GetLabels())) {
				ret = append(ret, obj)
			}
		}
	}
	return ret, err
}

// Get retrieves the ReplicationController from the indexer for a given workspace, namespace and name.
func (s *replicationControllerNamespaceLister) Get(name string) (*corev1.ReplicationController, error) {
	key := kcpcache.ToClusterAwareKey(s.cluster.String(), s.namespace, name)
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(corev1.Resource("ReplicationController"), name)
	}
	return obj.(*corev1.ReplicationController), nil
}
