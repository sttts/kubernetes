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
	json "encoding/json"
	"fmt"
	"time"

	v1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	admissionregistrationv1 "k8s.io/client-go/applyconfigurations/admissionregistration/v1"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// MutatingWebhookConfigurationsGetter has a method to return a MutatingWebhookConfigurationInterface.
// A group's client should implement this interface.
type MutatingWebhookConfigurationsGetter interface {
	MutatingWebhookConfigurations() MutatingWebhookConfigurationInterface
}

type ScopedMutatingWebhookConfigurationsGetter interface {
	ScopedMutatingWebhookConfigurations(scope rest.Scope) MutatingWebhookConfigurationInterface
}

// MutatingWebhookConfigurationInterface has methods to work with MutatingWebhookConfiguration resources.
type MutatingWebhookConfigurationInterface interface {
	Create(ctx context.Context, mutatingWebhookConfiguration *v1.MutatingWebhookConfiguration, opts metav1.CreateOptions) (*v1.MutatingWebhookConfiguration, error)
	Update(ctx context.Context, mutatingWebhookConfiguration *v1.MutatingWebhookConfiguration, opts metav1.UpdateOptions) (*v1.MutatingWebhookConfiguration, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.MutatingWebhookConfiguration, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.MutatingWebhookConfigurationList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.MutatingWebhookConfiguration, err error)
	Apply(ctx context.Context, mutatingWebhookConfiguration *admissionregistrationv1.MutatingWebhookConfigurationApplyConfiguration, opts metav1.ApplyOptions) (result *v1.MutatingWebhookConfiguration, err error)
	MutatingWebhookConfigurationExpansion
}

// mutatingWebhookConfigurations implements MutatingWebhookConfigurationInterface
type mutatingWebhookConfigurations struct {
	client  rest.Interface
	cluster string
	scope   rest.Scope
}

// newMutatingWebhookConfigurations returns a MutatingWebhookConfigurations
func newMutatingWebhookConfigurations(c *AdmissionregistrationV1Client, scope rest.Scope) *mutatingWebhookConfigurations {
	return &mutatingWebhookConfigurations{
		client:  c.RESTClient(),
		cluster: c.cluster,
		scope:   scope,
	}
}

// Get takes name of the mutatingWebhookConfiguration, and returns the corresponding mutatingWebhookConfiguration object, and an error if there is any.
func (c *mutatingWebhookConfigurations) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.MutatingWebhookConfiguration, err error) {
	result = &v1.MutatingWebhookConfiguration{}
	err = c.client.Get().
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MutatingWebhookConfigurations that match those selectors.
func (c *mutatingWebhookConfigurations) List(ctx context.Context, opts metav1.ListOptions) (result *v1.MutatingWebhookConfigurationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.MutatingWebhookConfigurationList{}
	err = c.client.Get().
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mutatingWebhookConfigurations.
func (c *mutatingWebhookConfigurations) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a mutatingWebhookConfiguration and creates it.  Returns the server's representation of the mutatingWebhookConfiguration, and an error, if there is any.
func (c *mutatingWebhookConfigurations) Create(ctx context.Context, mutatingWebhookConfiguration *v1.MutatingWebhookConfiguration, opts metav1.CreateOptions) (result *v1.MutatingWebhookConfiguration, err error) {
	result = &v1.MutatingWebhookConfiguration{}
	err = c.client.Post().
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(mutatingWebhookConfiguration).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a mutatingWebhookConfiguration and updates it. Returns the server's representation of the mutatingWebhookConfiguration, and an error, if there is any.
func (c *mutatingWebhookConfigurations) Update(ctx context.Context, mutatingWebhookConfiguration *v1.MutatingWebhookConfiguration, opts metav1.UpdateOptions) (result *v1.MutatingWebhookConfiguration, err error) {
	result = &v1.MutatingWebhookConfiguration{}
	err = c.client.Put().
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		Name(mutatingWebhookConfiguration.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(mutatingWebhookConfiguration).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the mutatingWebhookConfiguration and deletes it. Returns an error if one occurs.
func (c *mutatingWebhookConfigurations) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mutatingWebhookConfigurations) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched mutatingWebhookConfiguration.
func (c *mutatingWebhookConfigurations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.MutatingWebhookConfiguration, err error) {
	result = &v1.MutatingWebhookConfiguration{}
	err = c.client.Patch(pt).
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied mutatingWebhookConfiguration.
func (c *mutatingWebhookConfigurations) Apply(ctx context.Context, mutatingWebhookConfiguration *admissionregistrationv1.MutatingWebhookConfigurationApplyConfiguration, opts metav1.ApplyOptions) (result *v1.MutatingWebhookConfiguration, err error) {
	if mutatingWebhookConfiguration == nil {
		return nil, fmt.Errorf("mutatingWebhookConfiguration provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(mutatingWebhookConfiguration)
	if err != nil {
		return nil, err
	}
	name := mutatingWebhookConfiguration.Name
	if name == nil {
		return nil, fmt.Errorf("mutatingWebhookConfiguration.Name must be provided to Apply")
	}
	result = &v1.MutatingWebhookConfiguration{}
	err = c.client.Patch(types.ApplyPatchType).
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("mutatingwebhookconfigurations").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
