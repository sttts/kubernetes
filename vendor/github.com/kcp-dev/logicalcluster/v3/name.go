/*
Copyright 2022 The KCP Authors.

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

package logicalcluster

import "regexp"

var (
	// Wildcard is the Name indicating cross-workspace requests.
	Wildcard = Name{"*"}

	clusterNameRegExp = regexp.MustCompile(clusterNameString)
)

const (
	clusterNameString string = "^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?"
)

// Name represent a cluster name. A cluster name
// 1. can be used to access a cluster via `/cluster/<cluster-name>`
// 2. is part of an etcd key path
// 3. is used to uniquely reference a logical cluster. There is at most one cluster name for a logical cluster, but many logical cluster string can point to the same cluster name.
type Name struct {
	value string
}

// NewName returns a Name from a string.
func NewName(value string) Name {
	return Name{value}
}

// Path returns a Path (a logical cluster) out of the Name
func (n Name) Path() Path {
	return New(n.value)
}

// String returns the string representation of the cluster name.
func (n Name) String() string {
	return n.value
}

// IsValid returns true if the name is a Wildcard starts with a lower-case letter and contains only lower-case letters, digits and hyphens.
func (n Name) IsValid() bool {
	return n.value == "*" || clusterNameRegExp.MatchString(n.value)
}

// Empty returns true if the cluster name is unset.
func (n Name) Empty() bool {
	return n.value == ""
}

// Object is a local interface representation of the Kubernetes metav1.Object, to avoid dependencies on
// k8s.io/apimachinery.
type Object interface {
	GetAnnotations() map[string]string
}

// AnnotationKey is the name of the annotation key used to denote an object's logical cluster.
const AnnotationKey = "kcp.dev/cluster"

// From returns the logical cluster name for obj.
func From(obj Object) Name {
	return Name{obj.GetAnnotations()[AnnotationKey]}
}
