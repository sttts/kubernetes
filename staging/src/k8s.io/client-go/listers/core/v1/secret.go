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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// SecretLister helps list Secrets.
// All objects returned here must be treated as read-only.
type SecretLister interface {
	Scoped(scope rest.Scope) SecretLister
	// List lists all Secrets in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Secret, err error)
	// Secrets returns an object that can list and get Secrets.
	Secrets(namespace string) SecretNamespaceLister
	SecretListerExpansion
}

// secretLister implements the SecretLister interface.
type secretLister struct {
	indexer cache.Indexer
	scope   rest.Scope
}

// NewSecretLister returns a new SecretLister.
func NewSecretLister(indexer cache.Indexer) SecretLister {
	return &secretLister{indexer: indexer}
}

func (s *secretLister) Scoped(scope rest.Scope) SecretLister {
	return &secretLister{
		indexer: s.indexer,
		scope:   scope,
	}
}

// List lists all Secrets in the indexer.
func (s *secretLister) List(selector labels.Selector) (ret []*v1.Secret, err error) {
	appendFunc := func(m interface{}) {
		ret = append(ret, m.(*v1.Secret))
	}

	if s.scope == nil {
		err = cache.ListAll(s.indexer, selector, appendFunc)
		return ret, err
	}

	indexValue := s.scope.Name()
	err = cache.ListAllByIndexAndValue(s.indexer, cache.ListAllIndex, indexValue, selector, appendFunc)
	return ret, err
}

// Secrets returns an object that can list and get Secrets.
func (s *secretLister) Secrets(namespace string) SecretNamespaceLister {
	return secretNamespaceLister{indexer: s.indexer, namespace: namespace, scope: s.scope}
}

// SecretNamespaceLister helps list and get Secrets.
// All objects returned here must be treated as read-only.
type SecretNamespaceLister interface {
	// List lists all Secrets in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Secret, err error)
	// Get retrieves the Secret from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Secret, error)
	SecretNamespaceListerExpansion
}

// secretNamespaceLister implements the SecretNamespaceLister
// interface.
type secretNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
	scope     rest.Scope
}

// List lists all Secrets in the indexer for a given namespace.
func (s secretNamespaceLister) List(selector labels.Selector) (ret []*v1.Secret, err error) {
	indexValue := s.namespace
	if s.scope != nil {
		indexValue = s.scope.CacheKey(s.namespace)
	}
	err = cache.ListAllByIndexAndValue(s.indexer, cache.NamespaceIndex, indexValue, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Secret))
	})
	return ret, err
}

// Get retrieves the Secret from the indexer for a given namespace and name.
func (s secretNamespaceLister) Get(name string) (*v1.Secret, error) {
	key := s.namespace + "/" + name
	if s.scope != nil {
		key = s.scope.CacheKey(key)
	}
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("secret"), name)
	}
	return obj.(*v1.Secret), nil
}
