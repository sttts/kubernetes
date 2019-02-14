package openapi

import (
	"context"
	"fmt"
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	generatedopenapi "k8s.io/apiextensions-apiserver/pkg/generated/openapi"
	"k8s.io/apiextensions-apiserver/pkg/registry/customresource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	kubeopenapibuilder "k8s.io/kube-openapi/pkg/builder"
	openapicommon "k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/util"
)

// newSpec creates an OpenAPI spec for the given crd, version and schema.
func newSpec(crd *apiextensions.CustomResourceDefinition, version string, openAPISchema *spec.Schema) (*spec.Swagger, error) {
	// TODO: do not mutate
	//openAPISchema.SetProperty("metadata", *spec.RefSchema(objectMetaSchemaRef).
	//	WithDescription(swaggerPartialObjectMetadataDescriptions["metadata"]))
	//openAPISchema.SetProperty("apiVersion", getDefinition(typeMetaType).SchemaProps.Properties["apiVersion"])
	//openAPISchema.SetProperty("kind", getDefinition(typeMetaType).SchemaProps.Properties["kind"])
	openAPISchema.AddExtension(endpoints.ROUTE_META_GVK, []map[string]string{
		{
			"group":   crd.Spec.Group,
			"version": version,
			"kind":    crd.Status.AcceptedNames.Kind,
		},
	})

	rootScoped := sets.NewString()
	if crd.Spec.Scope == apiextensions.ClusterScoped {
		rootScoped.Insert(crd.Status.AcceptedNames.Kind)
	}
	gvk := schema.GroupVersionKind{Group: crd.Spec.Group, Version: version, Kind: crd.Status.AcceptedNames.Kind}
	var status *apiextensions.CustomResourceSubresourceStatus
	var scale *apiextensions.CustomResourceSubresourceScale
	if crd.Spec.Subresources != nil {
		status = crd.Spec.Subresources.Status
		scale = crd.Spec.Subresources.Scale
	}
	scheme := runtime.NewScheme()
	metav1.AddToGroupVersion(scheme, gvk.GroupVersion())
	creater := dummyCreater{gvk, scheme}
	typer := dummyTyper{gvk}
	strategy := customresource.NewStrategy(typer, crd.Spec.Scope == apiextensions.NamespaceScoped, gvk, nil, nil, status, scale)
	storage := customresource.NewStorage(crd, version,
		strategy,
		generic.RESTOptions{Decorator: generic.UndecoratedStorage, ResourcePrefix: "/fake", StorageConfig: &storagebackend.Config{Type: "etcd3"}},
		nil,
		&dummyStorage{},
	)

	restStorages := map[string]rest.Storage{
		crd.Status.AcceptedNames.Plural: storage.CustomResource,
	}
	subresources, err := apiextensions.GetSubresourcesForVersion(crd, version)
	if err != nil {
		return nil, err
	}
	if subresources != nil && subresources.Status != nil {
		restStorages[crd.Status.AcceptedNames.Plural+"/status"] = storage.Status
	}
	if subresources != nil && subresources.Scale != nil {
		restStorages[crd.Status.AcceptedNames.Plural+"/scale"] = storage.Scale
	}

	apiGroupVersion := endpoints.APIGroupVersion{
		Storage:                      restStorages,
		Root:                         "/apis",
		GroupVersion:                 schema.GroupVersion{crd.Spec.Group, version},
		OptionsExternalVersion:       nil,
		MetaGroupVersion:             &metav1.SchemeGroupVersion,
		RootScopedKinds:              rootScoped,
		Serializer:                   jsonDummySerialize{creater, typer},
		ParameterCodec:               nil,
		Typer:                        typer,
		Creater:                      creater,
		Convertor:                    nil,
		Defaulter:                    nil,
		Linker:                       nil,
		UnsafeConvertor:              nil,
		Authorizer:                   nil,
		Admit:                        nil,
		MinRequestTimeout:            0,
		EnableAPIResponseCompression: false,
		OpenAPIModels:                nil,
	}
	container := restful.NewContainer()
	if err := apiGroupVersion.InstallREST(container); err != nil {
		return nil, err
	}

	config := crdOpenAPIConfig(crd, version, openAPISchema, openAPISchema /* TODODODODODOD */)
	spec, err := kubeopenapibuilder.BuildOpenAPISpec(container.RegisteredWebServices(), config)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// getOpenAPIConfig builds config which wires up generated definitions for kube-openapi to consume
func crdOpenAPIConfig(crd *apiextensions.CustomResourceDefinition, version string, kindSchema, listSchema *spec.Schema) *openapicommon.Config {
	namer := openapi.DefinitionNamer{}
	return &openapicommon.Config{
		ProtocolList: []string{"https"},
		Info: &spec.Info{
			InfoProps: spec.InfoProps{
				Title:   "Kubernetes CRD Swagger",
				Version: "v0.1.0",
			},
		},
		CommonResponses: map[int]spec.Response{
			401: {
				ResponseProps: spec.ResponseProps{
					Description: "Unauthorized",
				},
			},
		},
		GetOperationIDAndTags: openapi.GetOperationIDAndTags,
		GetDefinitionName:     namer.GetDefinitionName,
		GetDefinitions: func(rc openapicommon.ReferenceCallback) map[string]openapicommon.OpenAPIDefinition {
			// copy static definitions
			staticDefs := generatedopenapi.GetOpenAPIDefinitions(rc)
			defs := make(map[string]openapicommon.OpenAPIDefinition, len(staticDefs))
			for k, v := range staticDefs {
				defs[k] = v
			}

			// register CRD and its list variant
			gvk := schema.GroupVersionKind{Group: crd.Spec.Group, Version: version, Kind: crd.Status.AcceptedNames.Kind}
			listGVK := schema.GroupVersionKind{Group: crd.Spec.Group, Version: version, Kind: crd.Status.AcceptedNames.ListKind}
			defs[canonicalTypeName(gvk)] = openapicommon.OpenAPIDefinition{Schema: *kindSchema}
			defs[canonicalTypeName(listGVK)] = openapicommon.OpenAPIDefinition{Schema: *listSchema}
			defs["k8s.io/apimachinery/pkg/runtime.RawExtension"] = openapicommon.OpenAPIDefinition{}

			return defs
		},
	}
}

type dummyObject struct {
	dummy int
	gvk   schema.GroupVersionKind
}

var (
	_ runtime.Object          = &dummyObject{}
	_ util.CanonicalTypeNamer = dummyObject{}
)

func (o *dummyObject) DeepCopyObject() runtime.Object {
	return o
}

func (o *dummyObject) GetObjectKind() schema.ObjectKind {
	return o
}

func (o *dummyObject) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	o.gvk = gvk
}

func (o *dummyObject) GroupVersionKind() schema.GroupVersionKind {
	return o.gvk
}

func (o dummyObject) CanonicalTypeName() string {
	return canonicalTypeName(o.gvk)
}

func canonicalTypeName(gvk schema.GroupVersionKind) string {
	cs := strings.Split(gvk.Group, ".")
	if len(cs) == 1 {
		return fmt.Sprintf("%s/%s.%s", cs[0], gvk.Version, gvk.Kind)
	} else if len(cs) == 2 {
		return fmt.Sprintf("%s.%s/%s.%s", cs[0], cs[1], gvk.Version, gvk.Kind)
	}
	return fmt.Sprintf("%s.%s/%s/%s.%s", cs[0], cs[1], strings.Join(cs[2:], "/"), gvk.Version, gvk.Kind)
}

type dummyStorage struct {
}

func (d *dummyStorage) Versioner() storage.Versioner { return nil }
func (d *dummyStorage) Create(_ context.Context, _ string, _, _ runtime.Object, _ uint64) error {
	return fmt.Errorf("unimplemented")
}
func (d *dummyStorage) Delete(_ context.Context, _ string, _ runtime.Object, _ *storage.Preconditions) error {
	return fmt.Errorf("unimplemented")
}
func (d *dummyStorage) Watch(_ context.Context, _ string, _ string, _ storage.SelectionPredicate) (watch.Interface, error) {
	return nil, fmt.Errorf("unimplemented")
}
func (d *dummyStorage) WatchList(_ context.Context, _ string, _ string, _ storage.SelectionPredicate) (watch.Interface, error) {
	return nil, fmt.Errorf("unimplemented")
}
func (d *dummyStorage) Get(_ context.Context, _ string, _ string, _ runtime.Object, _ bool) error {
	return fmt.Errorf("unimplemented")
}
func (d *dummyStorage) GetToList(_ context.Context, _ string, _ string, _ storage.SelectionPredicate, _ runtime.Object) error {
	return fmt.Errorf("unimplemented")
}
func (d *dummyStorage) List(_ context.Context, _ string, _ string, _ storage.SelectionPredicate, listObj runtime.Object) error {
	return fmt.Errorf("unimplemented")
}
func (d *dummyStorage) GuaranteedUpdate(_ context.Context, _ string, _ runtime.Object, _ bool, _ *storage.Preconditions, _ storage.UpdateFunc, _ ...runtime.Object) error {
	return fmt.Errorf("unimplemented")
}
func (d *dummyStorage) Count(_ string) (int64, error) {
	return 0, fmt.Errorf("unimplemented")
}

type jsonDummySerialize struct {
	creater runtime.ObjectCreater
	typer   runtime.ObjectTyper
}

func (s jsonDummySerialize) SupportedMediaTypes() []runtime.SerializerInfo {
	return []runtime.SerializerInfo{
		{
			MediaType:        "application/json",
			EncodesAsText:    true,
			Serializer:       json.NewSerializer(json.DefaultMetaFactory, s.creater, s.typer, false),
			PrettySerializer: json.NewSerializer(json.DefaultMetaFactory, s.creater, s.typer, true),
			StreamSerializer: &runtime.StreamSerializerInfo{
				EncodesAsText: true,
				Serializer:    json.NewSerializer(json.DefaultMetaFactory, s.creater, s.typer, false),
				Framer:        json.Framer,
			},
		},
		{
			MediaType:     "application/yaml",
			EncodesAsText: true,
			Serializer:    json.NewYAMLSerializer(json.DefaultMetaFactory, s.creater, s.typer),
		},
	}
}

func (s jsonDummySerialize) EncoderForVersion(encoder runtime.Encoder, gv runtime.GroupVersioner) runtime.Encoder {
	return nil
}

func (s jsonDummySerialize) DecoderToVersion(decoder runtime.Decoder, gv runtime.GroupVersioner) runtime.Decoder {
	return nil
}

type dummyTyper struct {
	gvk schema.GroupVersionKind
}

func (t dummyTyper) ObjectKinds(runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	return []schema.GroupVersionKind{t.gvk}, false, nil
}

func (t dummyTyper) Recognizes(gvk schema.GroupVersionKind) bool {
	return t.gvk == gvk
}

type dummyCreater struct {
	gvk    schema.GroupVersionKind
	scheme *runtime.Scheme
}

func (c dummyCreater) New(kind schema.GroupVersionKind) (out runtime.Object, err error) {
	if kind == c.gvk {
		return &dummyObject{0, kind}, nil
	}
	return c.scheme.New(kind)
}
