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

package v1

import (
	v1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// EvictionLister helps list Evictions.
// All objects returned here must be treated as read-only.
type EvictionLister interface {
	Scoped(scope rest.Scope) EvictionLister
	// List lists all Evictions in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Eviction, err error)
	// Evictions returns an object that can list and get Evictions.
	Evictions(namespace string) EvictionNamespaceLister
	EvictionListerExpansion
}

// evictionLister implements the EvictionLister interface.
type evictionLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewEvictionLister returns a new EvictionLister.
func NewEvictionLister(indexer cache.Indexer) EvictionLister {
	return &evictionLister{indexer: indexer}
}

func (s *evictionLister) Scoped(scope rest.Scope) EvictionLister {
	return &evictionLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all Evictions in the indexer.
func (s *evictionLister) List(selector labels.Selector) (ret []*v1.Eviction, err error) {
	var indexValue string
	if s.scope != nil {
		indexValue = s.scope.Name()
	}
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Eviction))
	})
	return ret, err
}

// Evictions returns an object that can list and get Evictions.
func (s *evictionLister) Evictions(namespace string) EvictionNamespaceLister {
	return evictionNamespaceLister{indexer: s.indexer, namespace: namespace, scope: s.scope}
}

// EvictionNamespaceLister helps list and get Evictions.
// All objects returned here must be treated as read-only.
type EvictionNamespaceLister interface {
	// List lists all Evictions in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Eviction, err error)
	// Get retrieves the Eviction from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Eviction, error)
	EvictionNamespaceListerExpansion
}

// evictionNamespaceLister implements the EvictionNamespaceLister
// interface.
type evictionNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
	scope     rest.Scope
}

// List lists all Evictions in the indexer for a given namespace.
func (s evictionNamespaceLister) List(selector labels.Selector) (ret []*v1.Eviction, err error) {
	indexValue := s.namespace
	if s.scope != nil {
		indexValue = s.scope.CacheKey(s.namespace)
	}
	err = cache.ListAllByIndexAndValue(s.indexer, cache.NamespaceIndex, indexValue, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Eviction))
	})
	return ret, err
}

// Get retrieves the Eviction from the indexer for a given namespace and name.
func (s evictionNamespaceLister) Get(name string) (*v1.Eviction, error) {
	key := s.namespace + "/" + name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("eviction"), name)
	}
	return obj.(*v1.Eviction), nil
}
