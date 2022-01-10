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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1 "k8s.io/code-generator/examples/HyphenGroup/apis/example/v1"
	scheme "k8s.io/code-generator/examples/HyphenGroup/clientset/versioned/scheme"
)

// TestTypesGetter has a method to return a TestTypeInterface.
// A group's client should implement this interface.
type TestTypesGetter interface {
	TestTypes(namespace string) TestTypeInterface
}

type ScopedTestTypesGetter interface {
	ScopedTestTypes(scope rest.Scope, namespace string) TestTypeInterface
}

// TestTypeInterface has methods to work with TestType resources.
type TestTypeInterface interface {
	Create(ctx context.Context, testType *v1.TestType, opts metav1.CreateOptions) (*v1.TestType, error)
	Update(ctx context.Context, testType *v1.TestType, opts metav1.UpdateOptions) (*v1.TestType, error)
	UpdateStatus(ctx context.Context, testType *v1.TestType, opts metav1.UpdateOptions) (*v1.TestType, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.TestType, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.TestTypeList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.TestType, err error)
	TestTypeExpansion
}

// testTypes implements TestTypeInterface
type testTypes struct {
	client  rest.Interface
	cluster string
	scope   rest.Scope
	ns      string
}

// newTestTypes returns a TestTypes
func newTestTypes(c *ExampleGroupV1Client, scope rest.Scope, namespace string) *testTypes {
	return &testTypes{
		client:  c.RESTClient(),
		cluster: c.cluster,
		scope:   scope,
		ns:      namespace,
	}
}

// Get takes name of the testType, and returns the corresponding testType object, and an error if there is any.
func (c *testTypes) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.TestType, err error) {
	result = &v1.TestType{}
	err = c.client.Get().
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TestTypes that match those selectors.
func (c *testTypes) List(ctx context.Context, opts metav1.ListOptions) (result *v1.TestTypeList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.TestTypeList{}
	err = c.client.Get().
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested testTypes.
func (c *testTypes) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a testType and creates it.  Returns the server's representation of the testType, and an error, if there is any.
func (c *testTypes) Create(ctx context.Context, testType *v1.TestType, opts metav1.CreateOptions) (result *v1.TestType, err error) {
	result = &v1.TestType{}
	err = c.client.Post().
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(testType).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a testType and updates it. Returns the server's representation of the testType, and an error, if there is any.
func (c *testTypes) Update(ctx context.Context, testType *v1.TestType, opts metav1.UpdateOptions) (result *v1.TestType, err error) {
	result = &v1.TestType{}
	err = c.client.Put().
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		Name(testType.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(testType).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *testTypes) UpdateStatus(ctx context.Context, testType *v1.TestType, opts metav1.UpdateOptions) (result *v1.TestType, err error) {
	result = &v1.TestType{}
	err = c.client.Put().
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		Name(testType.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(testType).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the testType and deletes it. Returns an error if one occurs.
func (c *testTypes) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *testTypes) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched testType.
func (c *testTypes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.TestType, err error) {
	result = &v1.TestType{}
	err = c.client.Patch(pt).
		Cluster(c.cluster).
		Scope(c.scope).
		Namespace(c.ns).
		Resource("testtypes").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
