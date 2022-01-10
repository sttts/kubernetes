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
	v1alpha1 "k8s.io/api/scheduling/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// PriorityClassLister helps list PriorityClasses.
// All objects returned here must be treated as read-only.
type PriorityClassLister interface {
	Scoped(scope rest.Scope) PriorityClassLister
	// List lists all PriorityClasses in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.PriorityClass, err error)
	// Get retrieves the PriorityClass from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.PriorityClass, error)
	PriorityClassListerExpansion
}

// priorityClassLister implements the PriorityClassLister interface.
type priorityClassLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewPriorityClassLister returns a new PriorityClassLister.
func NewPriorityClassLister(indexer cache.Indexer) PriorityClassLister {
	return &priorityClassLister{indexer: indexer}
}

func (s *priorityClassLister) Scoped(scope rest.Scope) PriorityClassLister {
	return &priorityClassLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all PriorityClasses in the indexer.
func (s *priorityClassLister) List(selector labels.Selector) (ret []*v1alpha1.PriorityClass, err error) {
	var indexValue string
	if s.scope != nil {
		indexValue = s.scope.Name()
	}
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.PriorityClass))
	})
	return ret, err
}

// Get retrieves the PriorityClass from the index for a given name.
func (s *priorityClassLister) Get(name string) (*v1alpha1.PriorityClass, error) {
	key := name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("priorityclass"), name)
	}
	return obj.(*v1alpha1.PriorityClass), nil
}
