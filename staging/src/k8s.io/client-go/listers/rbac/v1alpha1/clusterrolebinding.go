/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"

	v1alpha1 "k8s.io/api/rbac/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ClusterRoleBindingLister helps list ClusterRoleBindings.
// All objects returned here must be treated as read-only.
type ClusterRoleBindingLister interface {
	// List lists all ClusterRoleBindings in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ClusterRoleBinding, err error)
	// ListWithContext lists all ClusterRoleBindings in the indexer.
	// Objects returned here must be treated as read-only.
	ListWithContext(ctx context.Context, selector labels.Selector) (ret []*v1alpha1.ClusterRoleBinding, err error)
	// Get retrieves the ClusterRoleBinding from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ClusterRoleBinding, error)
	// GetWithContext retrieves the ClusterRoleBinding from the index for a given name.
	// Objects returned here must be treated as read-only.
	GetWithContext(ctx context.Context, name string) (*v1alpha1.ClusterRoleBinding, error)
	ClusterRoleBindingListerExpansion
}

// clusterRoleBindingLister implements the ClusterRoleBindingLister interface.
type clusterRoleBindingLister struct {
	indexer cache.Indexer
}

// NewClusterRoleBindingLister returns a new ClusterRoleBindingLister.
func NewClusterRoleBindingLister(indexer cache.Indexer) ClusterRoleBindingLister {
	return &clusterRoleBindingLister{indexer: indexer}
}

// List lists all ClusterRoleBindings in the indexer.
func (s *clusterRoleBindingLister) List(selector labels.Selector) (ret []*v1alpha1.ClusterRoleBinding, err error) {
	return s.ListWithContext(context.Background(), selector)
}

// ListWithContext lists all ClusterRoleBindings in the indexer.
func (s *clusterRoleBindingLister) ListWithContext(ctx context.Context, selector labels.Selector) (ret []*v1alpha1.ClusterRoleBinding, err error) {
	err = cache.IndexedListAll(ctx, s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ClusterRoleBinding))
	})
	return ret, err
}

// Get retrieves the ClusterRoleBinding from the index for a given name.
func (s *clusterRoleBindingLister) Get(name string) (*v1alpha1.ClusterRoleBinding, error) {
	return s.GetWithContext(context.Background(), name)
}

// GetWithContext retrieves the ClusterRoleBinding from the index for a given name.
func (s *clusterRoleBindingLister) GetWithContext(ctx context.Context, name string) (*v1alpha1.ClusterRoleBinding, error) {
	key, err := cache.NameKeyFunc(ctx, name)
	if err != nil {
		return nil, err
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("clusterrolebinding"), name)
	}
	return obj.(*v1alpha1.ClusterRoleBinding), nil
}
