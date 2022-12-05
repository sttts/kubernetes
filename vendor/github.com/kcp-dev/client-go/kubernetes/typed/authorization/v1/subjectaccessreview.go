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
	kcpclient "github.com/kcp-dev/apimachinery/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	authorizationv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"
)

// SubjectAccessReviewsClusterGetter has a method to return a SubjectAccessReviewClusterInterface.
// A group's cluster client should implement this interface.
type SubjectAccessReviewsClusterGetter interface {
	SubjectAccessReviews() SubjectAccessReviewClusterInterface
}

// SubjectAccessReviewClusterInterface can scope down to one cluster and return a authorizationv1client.SubjectAccessReviewInterface.
type SubjectAccessReviewClusterInterface interface {
	Cluster(logicalcluster.Path) authorizationv1client.SubjectAccessReviewInterface
}

type subjectAccessReviewsClusterInterface struct {
	clientCache kcpclient.Cache[*authorizationv1client.AuthorizationV1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *subjectAccessReviewsClusterInterface) Cluster(name logicalcluster.Path) authorizationv1client.SubjectAccessReviewInterface {
	if name == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return c.clientCache.ClusterOrDie(name).SubjectAccessReviews()
}
