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

package v1beta1

import (
	"context"

	v1beta1 "k8s.io/api/authorization/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// SubjectAccessReviewsGetter has a method to return a SubjectAccessReviewInterface.
// A group's client should implement this interface.
type SubjectAccessReviewsGetter interface {
	SubjectAccessReviews() SubjectAccessReviewInterface
}

type ScopedSubjectAccessReviewsGetter interface {
	ScopedSubjectAccessReviews(scope rest.Scope) SubjectAccessReviewInterface
}

// SubjectAccessReviewInterface has methods to work with SubjectAccessReview resources.
type SubjectAccessReviewInterface interface {
	Create(ctx context.Context, subjectAccessReview *v1beta1.SubjectAccessReview, opts v1.CreateOptions) (*v1beta1.SubjectAccessReview, error)
	SubjectAccessReviewExpansion
}

// subjectAccessReviews implements SubjectAccessReviewInterface
type subjectAccessReviews struct {
	client  rest.Interface
	cluster string
	scope   rest.Scope
}

// newSubjectAccessReviews returns a SubjectAccessReviews
func newSubjectAccessReviews(c *AuthorizationV1beta1Client, scope rest.Scope) *subjectAccessReviews {
	return &subjectAccessReviews{
		client:  c.RESTClient(),
		cluster: c.cluster,
		scope:   scope,
	}
}

// Create takes the representation of a subjectAccessReview and creates it.  Returns the server's representation of the subjectAccessReview, and an error, if there is any.
func (c *subjectAccessReviews) Create(ctx context.Context, subjectAccessReview *v1beta1.SubjectAccessReview, opts v1.CreateOptions) (result *v1beta1.SubjectAccessReview, err error) {
	result = &v1beta1.SubjectAccessReview{}
	err = c.client.Post().
		Cluster(c.cluster).
		Scope(c.scope).
		Resource("subjectaccessreviews").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(subjectAccessReview).
		Do(ctx).
		Into(result)
	return
}
