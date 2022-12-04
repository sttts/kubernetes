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

import (
	"encoding/json"
	"path"
	"regexp"
	"strings"
)

var (
	// Wildcard is the Path indicating cross-workspace requests.
	Wildcard = Path{value: "*"}

	// None is the name indicating a cluster-unaware context.
	None = New("")

	// TODO is a value created by automated refactoring tools that should be replaced by a real Name.
	TODO = None
)

const (

	// ClusterHeader set to "<lcluster>" on a request is an alternative to accessing the
	// cluster via /clusters/<lcluster>. With that the <lcluster> can be access via normal kube-like
	// /api and /apis endpoints.
	ClusterHeader = "X-Kubernetes-Cluster"

	separator = ":"
)

// Path represents a colon separated list of words describing a path in a logical cluster hierarchy, like a file path in a file-system.
type Path struct {
	value string
}

// New returns a Path from a string.
func New(value string) Path {
	return Path{value}
}

// NewValidated returns a Path from a string and whether it is a valid logical cluster.
// A valid logical cluster returns true on IsValid().
func NewValidated(value string) (Path, bool) {
	n := Path{value}
	return n, n.IsValid()
}

// Empty returns true if the logical cluster value is unset.
func (n Path) Empty() bool {
	return n.value == ""
}

// Name translates a logical cluster (a Path) into a Name.
// If the Path has a parent we cannot translate into a valid Name.
// Otherwise, the Path becomes a Name.
func (n Path) Name() (Name, bool) {
	if _, hasParent := n.Parent(); hasParent {
		return "", false
	}
	return Name(n.value), true
}

// RequestPath returns a path segment for the logical cluster to access its API.
func (n Path) RequestPath() string {
	return path.Join("/clusters", n.value)
}

// String returns the string representation of the logical cluster name.
func (n Path) String() string {
	return n.value
}

// Parent returns the parent logical cluster name of the given logical cluster name.
func (n Path) Parent() (Path, bool) {
	parent, _ := n.Split()
	return parent, parent.value != ""
}

// Split splits logical cluster immediately following the final colon,
// separating it into a parent logical cluster and name component.
// If there is no colon in path, Split returns an empty logical cluster name
// and name set to path.
func (n Path) Split() (parent Path, name string) {
	i := strings.LastIndex(n.value, separator)
	if i < 0 {
		return Path{}, n.value
	}
	return Path{n.value[:i]}, n.value[i+1:]
}

// Base returns the last component of the logical cluster name.
func (n Path) Base() string {
	_, name := n.Split()
	return name
}

// Join joins a parent logical cluster name and a name component.
func (n Path) Join(name string) Path {
	if n.value == "" {
		return Path{name}
	}
	return Path{n.value + separator + name}
}

func (n Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(&n.value)
}

func (n *Path) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	n.value = s
	return nil
}

func (n Path) HasPrefix(other Path) bool {
	return strings.HasPrefix(n.value, other.value)
}

const lclusterNameFmt string = "[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?"

var lclusterRegExp = regexp.MustCompile("^" + lclusterNameFmt + "(:" + lclusterNameFmt + ")*$")

// IsValid returns true if the name is a Wildcard or a colon separated list of words where each word
// starts with a lower-case letter and contains only lower-case letters, digits and hyphens.
func (n Path) IsValid() bool {
	return lclusterRegExp.MatchString(n.value)
}
