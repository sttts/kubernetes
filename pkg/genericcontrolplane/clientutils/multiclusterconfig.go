/*
Copyright 2014 The Kubernetes Authors.

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

// Package app does all of the work necessary to create a Kubernetes
// APIServer by binding together the API, master and APIServer infrastructure.
// It can be configured and called directly or via the hyperkube framework.
package clientutils

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/server"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/rest"
	_ "k8s.io/component-base/metrics/prometheus/workqueue" // for workqueue metric registration
	"k8s.io/klog/v2"
)

type multiClusterClientConfigRoundTripper struct {
	rt                  http.RoundTripper
	requestInfoResolver func() genericapirequest.RequestInfoResolver
	enabledOn           sets.String
	disableSharding     bool
}

// EnableMultiCluster allows uses a rountripper to hack the rest.Config used by
// client-go APIs, in order to add support for logical clusters for a given list of resources.
// - By default it enables "wildcard" cluster in list and watch request to enable searchng in all logical clusters
// - For other types of requests, tries to guess the logical cluster where the operation should occur
// - It finally sets the "X-Kubernetes-Cluster" header accordingly.
//
// This is a temporary hack and should be replaced by thoughtful and real support of logical clusters
// in the client-go layer
func EnableMultiCluster(config *rest.Config, apiServerConfig *genericapiserver.Config, disableSharding bool, enabledOnResources ...string) {
	config.ContentConfig.ContentType = "application/json"

	config.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		defaultResolver := defaultRequestInfoResolver()
		return &multiClusterClientConfigRoundTripper{
			rt: rt,
			requestInfoResolver: func() genericapirequest.RequestInfoResolver {
				if apiServerConfig != nil {
					return apiServerConfig.RequestInfoResolver
				} else {
					return defaultResolver
				}
			},
			enabledOn:       sets.NewString(enabledOnResources...),
			disableSharding: disableSharding,
		}
	})
}

func defaultRequestInfoResolver() genericapirequest.RequestInfoResolver {
	apiPrefixes := sets.NewString(strings.Trim(server.APIGroupPrefix, "/")) // all possible API prefixes
	legacyAPIPrefixes := sets.String{}                                      // APIPrefixes that won't have groups (legacy)
	apiPrefixes.Insert(strings.Trim(server.DefaultLegacyAPIPrefix, "/"))
	legacyAPIPrefixes.Insert(strings.Trim(server.DefaultLegacyAPIPrefix, "/"))

	return &request.RequestInfoFactory{
		APIPrefixes:          apiPrefixes,
		GrouplessAPIPrefixes: legacyAPIPrefixes,
	}
}

func (mcrt *multiClusterClientConfigRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = utilnet.CloneRequest(req)
	requestInfo, err := mcrt.requestInfoResolver().NewRequestInfo(req)
	if err != nil {
		return nil, err
	}
	// contextScope := rest.ScopeFrom(req.Context())
	// contextCluster := genericapirequest.ClusterFrom(req.Context())
	var contextCluster *genericapirequest.Cluster
	if requestInfo != nil &&
		mcrt.enabledOn.Has(requestInfo.Resource) {
		resourceClusterName := ""
		headerCluster := ""
		switch requestInfo.Verb {
		case "list", "watch":
			if contextCluster != nil && !contextCluster.Wildcard {
				headerCluster = contextCluster.Name
			} else {
				headerCluster = "*"
			}
		case "create", "update":
			err := func() error {
				// We dn't try to mutate the object here. Just guessing the ClusterName field from the body,
				// in order to set the header accordingly
				reader, err := req.GetBody()
				if err != nil {
					klog.Infof("DEBUG: Error when trying to read the request body from the multicluster client config roundtripper: %v", err)
					return err
				}
				defer reader.Close()
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(reader)
				if err != nil {
					klog.Infof("DEBUG: Error when trying to read the request body from the multicluster client config roundtripper: %v", err)
					return err
				}
				bytes := buf.Bytes()
				obj, _, err := unstructured.UnstructuredJSONScheme.Decode(bytes, nil, nil)
				if err != nil {
					klog.Infof("DEBUG: Error when trying to read the request body from the multicluster client config roundtripper: %v", err)
					return err
				}
				if s, ok := obj.(metav1.Object); ok {
					resourceClusterName = s.GetClusterName()
				}
				return nil
			}()
			if err != nil {
				return nil, err
			}
			fallthrough
		default:
			if resourceClusterName != "" {
				if contextCluster != nil && contextCluster.Name != resourceClusterName {
					return nil, errors.New("Resource cluster name " + resourceClusterName + " incompatible with context cluster name " + contextCluster.Name)
				}
				headerCluster = resourceClusterName
			} else {
				if contextCluster != nil { //
					if contextCluster.Wildcard {
						return nil, errors.New("Cluster should be set for request " + requestInfo.Verb + ", but instead wildcards were provided.")
					}
					headerCluster = contextCluster.Name
				}
			}
		}
		if headerCluster == "" {
			return nil, fmt.Errorf("Cluster should not be empty for request '%s' on resource '%s' (%s)", requestInfo.Verb, requestInfo.Resource, requestInfo.Path)
		}
		req.Header.Add("X-Kubernetes-Cluster", headerCluster)
	} else {
		if contextCluster != nil && contextCluster.Name != "" {
			req.Header.Add("X-Kubernetes-Cluster", contextCluster.Name)
		}
	}
	if mcrt.disableSharding {
		// let internal clients opt out of sharding for now
		req.Header.Add("X-Kubernetes-Sharded-Request", "true")
	}
	return mcrt.rt.RoundTrip(req)
}
