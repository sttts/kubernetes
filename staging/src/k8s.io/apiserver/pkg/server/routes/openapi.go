/*
Copyright 2016 The Kubernetes Authors.

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

package routes

import (
	"net/http"

	restful "github.com/emicklei/go-restful"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"k8s.io/apiserver/pkg/server/mux"
	builder2 "k8s.io/kube-openapi/pkg/builder"
	"k8s.io/kube-openapi/pkg/builder3"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/handler"
	"k8s.io/kube-openapi/pkg/handler3"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

// OpenAPI installs spec endpoints for each web service.
type OpenAPI struct {
	Config *common.Config
}

// OpenAPIServiceProvider is a hacky way to
// replace a single OpenAPIService by a provider which will
// provide an distinct openAPIService per scope.
// This is required to implement CRD tenancy and have the openAPI
// models be conistent with the current scope.
//
// However this is just a first step, since a better way
// would be to completly avoid the need of registering a OpenAPIService
// for each scope. See the addition comments below.
type OpenAPIServiceProvider interface {
	ForScope(scopeName string) *handler.OpenAPIService
	AddScope(scopeName string)
	RemoveScope(scopeName string)
	UpdateSpec(openapiSpec *spec.Swagger) error
}

type scopeAwarePathHandler struct {
	scopeName          string
	addHandlerForScope func(scopeName string, handler http.Handler)
}

func (c *scopeAwarePathHandler) Handle(path string, handler http.Handler) {
	c.addHandlerForScope(c.scopeName, handler)
}

// HACK: This is the implementation of OpenAPIServiceProvider
// that allows supporting several scopes for CRD tenancy.
//
// However this should be conisdered a temporary step, to cope with the
// current design of OpenAPI publishing. But having to register every scope
// creates more cost on creating scopes.
// Instead, we'd expect us to slowly refactor the openapi generation code so
// that it can be used dynamically, and time limited or size limited openapi caches
// would be used to serve the calculated version.
// Finally a development princple for the logical cluster prototype would be
// - don't do static registration of logical clusters
// - do lazy instantiation wherever possible so that starting a new logical cluster remains as cheap as possible
type openAPIServiceProvider struct {
	staticSpec                   *spec.Swagger
	defaultOpenAPIServiceHandler http.Handler
	defaultOpenAPIService        *handler.OpenAPIService
	openAPIServices              map[string]*handler.OpenAPIService
	handlers                     map[string]http.Handler
	path                         string
	mux                          *mux.PathRecorderMux
}

var _ OpenAPIServiceProvider = (*openAPIServiceProvider)(nil)

func (p *openAPIServiceProvider) ForScope(scopeName string) *handler.OpenAPIService {
	return p.openAPIServices[scopeName]
}

func (p *openAPIServiceProvider) AddScope(scopeName string) {
	if _, found := p.openAPIServices[scopeName]; !found {
		openAPIVersionedService, err := handler.NewOpenAPIService(p.staticSpec)
		if err != nil {
			klog.Fatalf("Failed to create OpenAPIService: %v", err)
		}

		if err = openAPIVersionedService.RegisterOpenAPIVersionedService(p.path, &scopeAwarePathHandler{
			scopeName: scopeName,
			addHandlerForScope: func(scopeName string, handler http.Handler) {
				p.handlers[scopeName] = handler
			},
		}); err != nil {
			klog.Fatalf("Failed to register versioned open api spec for root: %v", err)
		}
		p.openAPIServices[scopeName] = openAPIVersionedService
	}
}

func (p *openAPIServiceProvider) RemoveScope(scopeName string) {
	delete(p.openAPIServices, scopeName)
	delete(p.handlers, scopeName)
}

func (p *openAPIServiceProvider) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	scope := rest.ScopeFrom(req.Context())
	if scope == nil {
		p.defaultOpenAPIServiceHandler.ServeHTTP(resp, req)
		return
	}
	handler, found := p.handlers[scope.Name()]
	if !found {
		resp.WriteHeader(404)
		return
	}
	handler.ServeHTTP(resp, req)
}

func (o *openAPIServiceProvider) UpdateSpec(openapiSpec *spec.Swagger) (err error) {
	return o.defaultOpenAPIService.UpdateSpec(openapiSpec)
}

func (p *openAPIServiceProvider) Register() {
	defaultOpenAPIService, err := handler.NewOpenAPIService(p.staticSpec)
	if err != nil {
		klog.Fatalf("Failed to create OpenAPIService: %v", err)
	}

	err = defaultOpenAPIService.RegisterOpenAPIVersionedService(p.path, &scopeAwarePathHandler{
		scopeName: "",
		addHandlerForScope: func(scopeName string, handler http.Handler) {
			p.defaultOpenAPIServiceHandler = handler
		},
	})
	if err != nil {
		klog.Fatalf("Failed to register versioned open api spec for root: %v", err)
	}

	p.defaultOpenAPIService = defaultOpenAPIService
	p.mux.Handle(p.path, p)
}

// Install adds the SwaggerUI webservice to the given mux.
func (oa OpenAPI) InstallV2(c *restful.Container, mux *mux.PathRecorderMux) (OpenAPIServiceProvider, *spec.Swagger) {
	spec, err := builder2.BuildOpenAPISpec(c.RegisteredWebServices(), oa.Config)
	if err != nil {
		klog.Fatalf("Failed to build open api spec for root: %v", err)
	}
	spec.Definitions = handler.PruneDefaults(spec.Definitions)

	provider := &openAPIServiceProvider{
		mux:             mux,
		staticSpec:      spec,
		openAPIServices: map[string]*handler.OpenAPIService{},
		handlers:        map[string]http.Handler{},
		path:            "/openapi/v2",
	}

	provider.Register()

	return provider, spec
}

// InstallV3 adds the static group/versions defined in the RegisteredWebServices to the OpenAPI v3 spec
func (oa OpenAPI) InstallV3(c *restful.Container, mux *mux.PathRecorderMux) *handler3.OpenAPIService {
	openAPIVersionedService, err := handler3.NewOpenAPIService(nil)
	if err != nil {
		klog.Fatalf("Failed to create OpenAPIService: %v", err)
	}

	err = openAPIVersionedService.RegisterOpenAPIV3VersionedService("/openapi/v3", mux)
	if err != nil {
		klog.Fatalf("Failed to register versioned open api spec for root: %v", err)
	}

	grouped := make(map[string][]*restful.WebService)

	for _, t := range c.RegisteredWebServices() {
		// Strip the "/" prefix from the name
		gvName := t.RootPath()[1:]
		grouped[gvName] = []*restful.WebService{t}
	}

	for gv, ws := range grouped {
		spec, err := builder3.BuildOpenAPISpec(ws, oa.Config)
		if err != nil {
			klog.Errorf("Failed to build OpenAPI v3 for group %s, %q", gv, err)

		}
		openAPIVersionedService.UpdateGroupVersion(gv, spec)
	}
	return openAPIVersionedService
}
