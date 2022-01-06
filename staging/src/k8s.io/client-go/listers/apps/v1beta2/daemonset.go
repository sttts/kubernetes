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

package v1beta2

import (
	"context"

	v1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DaemonSetLister helps list DaemonSets.
// All objects returned here must be treated as read-only.
type DaemonSetLister interface {
	// List lists all DaemonSets in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.DaemonSet, err error)
	// ListWithContext lists all DaemonSets in the indexer.
	// Objects returned here must be treated as read-only.
	ListWithContext(ctx context.Context, selector labels.Selector) (ret []*v1beta2.DaemonSet, err error)
	// DaemonSets returns an object that can list and get DaemonSets.
	DaemonSets(namespace string) DaemonSetNamespaceLister
	DaemonSetListerExpansion
}

// daemonSetLister implements the DaemonSetLister interface.
type daemonSetLister struct {
	indexer cache.Indexer
}

// NewDaemonSetLister returns a new DaemonSetLister.
func NewDaemonSetLister(indexer cache.Indexer) DaemonSetLister {
	return &daemonSetLister{indexer: indexer}
}

// List lists all DaemonSets in the indexer.
func (s *daemonSetLister) List(selector labels.Selector) (ret []*v1beta2.DaemonSet, err error) {
	return s.ListWithContext(context.Background(), selector)
}

// ListWithContext lists all DaemonSets in the indexer.
func (s *daemonSetLister) ListWithContext(ctx context.Context, selector labels.Selector) (ret []*v1beta2.DaemonSet, err error) {
	err = cache.IndexedListAll(ctx, s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.DaemonSet))
	})
	return ret, err
}

// DaemonSets returns an object that can list and get DaemonSets.
func (s *daemonSetLister) DaemonSets(namespace string) DaemonSetNamespaceLister {
	return daemonSetNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DaemonSetNamespaceLister helps list and get DaemonSets.
// All objects returned here must be treated as read-only.
type DaemonSetNamespaceLister interface {
	// List lists all DaemonSets in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.DaemonSet, err error)
	// ListWithContext lists all DaemonSets in the indexer.
	// Objects returned here must be treated as read-only.
	ListWithContext(ctx context.Context, selector labels.Selector) (ret []*v1beta2.DaemonSet, err error)
	// Get retrieves the DaemonSet from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta2.DaemonSet, error)
	// GetWithContext retrieves the DaemonSet from the index for a given name.
	// Objects returned here must be treated as read-only.
	GetWithContext(ctx context.Context, name string) (*v1beta2.DaemonSet, error)
	DaemonSetNamespaceListerExpansion
}

// daemonSetNamespaceLister implements the DaemonSetNamespaceLister
// interface.
type daemonSetNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all DaemonSets in the indexer for a given namespace.
func (s daemonSetNamespaceLister) List(selector labels.Selector) (ret []*v1beta2.DaemonSet, err error) {
	return s.ListWithContext(context.Background(), selector)
}

// ListWithContext lists all DaemonSets in the indexer for a given namespace.
func (s daemonSetNamespaceLister) ListWithContext(ctx context.Context, selector labels.Selector) (ret []*v1beta2.DaemonSet, err error) {
	err = cache.ListAllByNamespace2(ctx, s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.DaemonSet))
	})
	return ret, err
}

// Get retrieves the DaemonSet from the indexer for a given namespace and name.
func (s daemonSetNamespaceLister) Get(name string) (*v1beta2.DaemonSet, error) {
	return s.GetWithContext(context.Background(), name)
}

// GetWithContext retrieves the DaemonSet from the indexer for a given namespace and name.
func (s daemonSetNamespaceLister) GetWithContext(ctx context.Context, name string) (*v1beta2.DaemonSet, error) {
	key, err := cache.NamespaceNameKeyFunc(ctx, s.namespace, name)
	if err != nil {
		return nil, err
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta2.Resource("daemonset"), name)
	}
	return obj.(*v1beta2.DaemonSet), nil
}
