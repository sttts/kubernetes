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
	v1alpha1 "k8s.io/api/rbac/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// RoleLister helps list Roles.
// All objects returned here must be treated as read-only.
type RoleLister interface {
	Scoped(scope rest.Scope) RoleLister
	// List lists all Roles in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Role, err error)
	// Roles returns an object that can list and get Roles.
	Roles(namespace string) RoleNamespaceLister
	RoleListerExpansion
}

// roleLister implements the RoleLister interface.
type roleLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewRoleLister returns a new RoleLister.
func NewRoleLister(indexer cache.Indexer) RoleLister {
	return &roleLister{indexer: indexer}
}

func (s *roleLister) Scoped(scope rest.Scope) RoleLister {
	return &roleLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all Roles in the indexer.
func (s *roleLister) List(selector labels.Selector) (ret []*v1alpha1.Role, err error) {
	appendFunc := func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Role))
	}

	if s.scope == nil {
		err = cache.ListAll(s.indexer, selector, appendFunc)
		return ret, err
	}

	indexValue := s.scope.Name()
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, appendFunc)
	return ret, err
}

// Roles returns an object that can list and get Roles.
func (s *roleLister) Roles(namespace string) RoleNamespaceLister {
	return roleNamespaceLister{indexer: s.indexer, namespace: namespace, scope: s.scope}
}

// RoleNamespaceLister helps list and get Roles.
// All objects returned here must be treated as read-only.
type RoleNamespaceLister interface {
	// List lists all Roles in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Role, err error)
	// Get retrieves the Role from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Role, error)
	RoleNamespaceListerExpansion
}

// roleNamespaceLister implements the RoleNamespaceLister
// interface.
type roleNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
	scope     rest.Scope
}

// List lists all Roles in the indexer for a given namespace.
func (s roleNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Role, err error) {
	indexValue := s.namespace
	if s.scope != nil {
		indexValue = s.scope.CacheKey(s.namespace)
	}
	err = cache.ListAllByIndexAndValue(s.indexer, cache.NamespaceIndex, indexValue, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Role))
	})
	return ret, err
}

// Get retrieves the Role from the indexer for a given namespace and name.
func (s roleNamespaceLister) Get(name string) (*v1alpha1.Role, error) {
	key := s.namespace + "/" + name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("role"), name)
	}
	return obj.(*v1alpha1.Role), nil
}
