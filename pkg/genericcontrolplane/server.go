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
package genericcontrolplane

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	apiextensionsapiserver "k8s.io/apiextensions-apiserver/pkg/apiserver"
	extensionsapiserver "k8s.io/apiextensions-apiserver/pkg/apiserver"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/authorization/union"
	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericfeatures "k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/filters"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/apiserver/pkg/util/feature"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	utilflowcontrol "k8s.io/apiserver/pkg/util/flowcontrol"
	"k8s.io/apiserver/pkg/util/notfoundhandler"
	"k8s.io/apiserver/pkg/util/openapi"
	"k8s.io/apiserver/pkg/util/webhook"
	clientgoinformers "k8s.io/client-go/informers"
	clientgoclientset "k8s.io/client-go/kubernetes"
	_ "k8s.io/component-base/metrics/prometheus/workqueue" // for workqueue metric registration
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/api/genericcontrolplanescheme"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	generatedopenapi "k8s.io/kubernetes/pkg/generated/openapi"
	"k8s.io/kubernetes/pkg/genericcontrolplane/aggregator"
	"k8s.io/kubernetes/pkg/genericcontrolplane/apis"
	"k8s.io/kubernetes/pkg/genericcontrolplane/clientutils"
	"k8s.io/kubernetes/pkg/genericcontrolplane/options"
	"k8s.io/kubernetes/pkg/kubeapiserver"
	"k8s.io/kubernetes/pkg/serviceaccount"
	"k8s.io/kubernetes/plugin/pkg/auth/authorizer/rbac"
)

const Include = "kube-control-plane"

const (
	etcdRetryLimit    = 60
	etcdRetryInterval = 1 * time.Second
	RootClusterName   = "admin"
)

// Run runs the specified APIServer.  This should never exit.
func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get())

	serverChain, err := CreateServerChain(completeOptions, stopCh)
	if err != nil {
		return err
	}
	server := serverChain.MiniAggregator.GenericAPIServer

	prepared := server.PrepareRun()
	return prepared.Run(stopCh)
}

type ServerChain struct {
	CustomResourceDefinitions *apiextensionsapiserver.CustomResourceDefinitions
	GenericControlPlane       *apis.GenericControlPlane
	MiniAggregator            *aggregator.MiniAggregatorServer
}

// unimplementedServiceResolver is a webhook.ServiceResolver that always returns an error, because
// we have not implemented support for this yet. As a result, CRD webhook conversions are not
// supported.
type unimplementedServiceResolver struct{}

// ResolveEndpoint always returns an error that this is not yet supported.
func (r *unimplementedServiceResolver) ResolveEndpoint(namespace string, name string, port int32) (*url.URL, error) {
	return nil, fmt.Errorf("CRD webhook conversions are not yet supported in kcp")
}

// CreateServerChain creates the apiservers connected via delegation.
func CreateServerChain(completedOptions completedServerRunOptions, stopCh <-chan struct{}) (*ServerChain, error) {
	kubeAPIServerConfig, pluginInitializer, err := CreateKubeAPIServerConfig(completedOptions)
	if err != nil {
		return nil, err
	}

	// Wire in a ServiceResolver that always returns an error that ResolveEndpoint is not yet
	// supported. The effect is that CRD webhook conversions are not supported and will always get an
	// error.
	serviceResolver := &unimplementedServiceResolver{}

	// If additional API servers are added, they should be gated.
	apiExtensionsConfig, err := createAPIExtensionsConfig(
		*kubeAPIServerConfig.GenericConfig,
		kubeAPIServerConfig.ExtraConfig.VersionedInformers,
		pluginInitializer,
		completedOptions.ServerRunOptions,
		serviceResolver,
		webhook.NewDefaultAuthenticationInfoResolverWrapper(
			nil,
			kubeAPIServerConfig.GenericConfig.EgressSelector,
			kubeAPIServerConfig.GenericConfig.LoopbackClientConfig,
			kubeAPIServerConfig.GenericConfig.TracerProvider,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("configure api extensions: %v", err)
	}
	notFoundHandler := notfoundhandler.New(kubeAPIServerConfig.GenericConfig.Serializer, genericapifilters.NoMuxAndDiscoveryIncompleteKey)
	apiExtensionsServer, err := createAPIExtensionsServer(apiExtensionsConfig, genericapiserver.NewEmptyDelegateWithCustomHandler(notFoundHandler))
	if err != nil {
		return nil, fmt.Errorf("create api extensions: %v", err)
	}

	kubeAPIServer, err := CreateKubeAPIServer(kubeAPIServerConfig, apiExtensionsServer.GenericAPIServer)
	if err != nil {
		return nil, err
	}

	miniAggregatorConfig := &aggregator.MiniAggregatorConfig{
		GenericConfig: kubeAPIServerConfig.GenericConfig,
	}

	if err := completedOptions.ServerRunOptions.Admission.ApplyTo(
		kubeAPIServerConfig.GenericConfig,
		kubeAPIServerConfig.ExtraConfig.VersionedInformers,
		kubeAPIServerConfig.GenericConfig.LoopbackClientConfig,
		feature.DefaultFeatureGate,
		pluginInitializer...); err != nil {
		return nil, err
	}

	miniAggregatorServer, err := miniAggregatorConfig.Complete(kubeAPIServerConfig.ExtraConfig.VersionedInformers).New(kubeAPIServer.GenericAPIServer, kubeAPIServer, apiExtensionsServer)
	if err != nil {
		return nil, err
	}

	return &ServerChain{
		CustomResourceDefinitions: apiExtensionsServer,
		GenericControlPlane:       kubeAPIServer,
		MiniAggregator:            miniAggregatorServer,
	}, nil
}

// CreateKubeAPIServer creates and wires a workable kube-apiserver
func CreateKubeAPIServer(kubeAPIServerConfig *apis.Config, delegateAPIServer genericapiserver.DelegationTarget) (*apis.GenericControlPlane, error) {
	kubeAPIServer, err := kubeAPIServerConfig.Complete().New(delegateAPIServer)
	if err != nil {
		return nil, err
	}

	return kubeAPIServer, nil
}

// CreateKubeAPIServerConfig creates all the resources for running the API server, but runs none of them
func CreateKubeAPIServerConfig(s completedServerRunOptions) (
	*apis.Config,
	[]admission.PluginInitializer,
	error,
) {
	genericConfig, pluginInitializers, storageFactory, err := BuildGenericConfig(s.ServerRunOptions)
	if err != nil {
		return nil, nil, err
	}

	s.Metrics.Apply()
	serviceaccount.RegisterMetrics()

	kubeClientConfig := genericConfig.LoopbackClientConfig
	clientgoExternalClient, err := clientgoclientset.NewForConfig(kubeClientConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create real external clientset: %v", err)
	}
	versionedInformers := clientgoinformers.NewSharedInformerFactory(clientgoExternalClient, 10*time.Minute)

	// TODO(ncdc,1.23) upstream has this in the Cobra RunE method and it's called as early as
	// possible. Do we want to consider doing something similar?
	if err := s.Logs.ValidateAndApply(); err != nil {
		return nil, nil, err
	}

	config := &apis.Config{
		GenericConfig: genericConfig,
		ExtraConfig: apis.ExtraConfig{
			APIResourceConfigSource: storageFactory.APIResourceConfigSource,
			StorageFactory:          storageFactory,
			EventTTL:                s.EventTTL,
			EnableLogsSupport:       s.EnableLogsHandler,

			VersionedInformers: versionedInformers,

			IdentityLeaseDurationSeconds:      s.IdentityLeaseDurationSeconds,
			IdentityLeaseRenewIntervalSeconds: s.IdentityLeaseRenewIntervalSeconds,
		},
	}

	clientCAProvider, err := s.Authentication.ClientCert.GetClientCAContentProvider()
	if err != nil {
		return nil, nil, err
	}
	config.ExtraConfig.ClusterAuthenticationInfo.ClientCA = clientCAProvider

	requestHeaderConfig, err := s.Authentication.RequestHeader.ToAuthenticationRequestHeaderConfig()
	if err != nil {
		return nil, nil, err
	}
	if requestHeaderConfig != nil {
		config.ExtraConfig.ClusterAuthenticationInfo.RequestHeaderCA = requestHeaderConfig.CAContentProvider
		config.ExtraConfig.ClusterAuthenticationInfo.RequestHeaderAllowedNames = requestHeaderConfig.AllowedClientNames
		config.ExtraConfig.ClusterAuthenticationInfo.RequestHeaderExtraHeaderPrefixes = requestHeaderConfig.ExtraHeaderPrefixes
		config.ExtraConfig.ClusterAuthenticationInfo.RequestHeaderGroupHeaders = requestHeaderConfig.GroupHeaders
		config.ExtraConfig.ClusterAuthenticationInfo.RequestHeaderUsernameHeaders = requestHeaderConfig.UsernameHeaders
	}

	// if err := config.GenericConfig.AddPostStartHook("start-kube-apiserver-admission-initializer", admissionPostStartHook); err != nil {
	// 	return nil, nil, nil, err
	// }

	// // Load the public keys.
	// var pubKeys []interface{}
	// for _, f := range s.Authentication.ServiceAccounts.KeyFiles {
	// 	keys, err := keyutil.PublicKeysFromFile(f)
	// 	if err != nil {
	// 		return nil, nil, nil, fmt.Errorf("failed to parse key file %q: %v", f, err)
	// 	}
	// 	pubKeys = append(pubKeys, keys...)
	// }
	// // Plumb the required metadata through ExtraConfig.
	// config.ExtraConfig.ServiceAccountIssuerURL = s.Authentication.ServiceAccounts.Issuers[0]
	// config.ExtraConfig.ServiceAccountJWKSURI = s.Authentication.ServiceAccounts.JWKSURI
	// config.ExtraConfig.ServiceAccountPublicKeys = pubKeys

	return config, pluginInitializers, nil
}

// BuildGenericConfig takes the master server options and produces the genericapiserver.Config associated with it
func BuildGenericConfig(
	s *options.ServerRunOptions,
) (
	genericConfig *genericapiserver.Config,
	pluginInitializers []admission.PluginInitializer,
	// admissionPostStartHook genericapiserver.PostStartHookFunc,
	storageFactory *serverstorage.DefaultStorageFactory,
	lastErr error,
) {
	genericConfig = genericapiserver.NewConfig(genericcontrolplanescheme.Codecs)
	if s.BuildHandlerChainFunc != nil {
		genericConfig.BuildHandlerChainFunc = s.BuildHandlerChainFunc
	}

	if lastErr = s.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = s.SecureServing.ApplyTo(&genericConfig.SecureServing, &genericConfig.LoopbackClientConfig); lastErr != nil {
		return
	}
	if lastErr = s.Features.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if lastErr = s.APIEnablement.ApplyTo(genericConfig, apis.DefaultAPIResourceConfigSource(), genericcontrolplanescheme.Scheme); lastErr != nil {
		return
	}
	if lastErr = s.EgressSelector.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	if utilfeature.DefaultFeatureGate.Enabled(genericfeatures.APIServerTracing) {
		if lastErr = s.Traces.ApplyTo(genericConfig.EgressSelector, genericConfig); lastErr != nil {
			return
		}
	}

	// wrap the definitions to revert any changes from disabled features
	getOpenAPIDefinitions := openapi.GetOpenAPIDefinitionsWithoutDisabledFeatures(generatedopenapi.GetOpenAPIDefinitions)
	genericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(getOpenAPIDefinitions, openapinamer.NewDefinitionNamer(legacyscheme.Scheme, extensionsapiserver.Scheme, extensionsapiserver.Scheme))
	genericConfig.OpenAPIConfig.Info.Title = "Kubernetes"
	genericConfig.LongRunningFunc = filters.BasicLongRunningRequestCheck(
		sets.NewString("watch", "proxy"),
		sets.NewString("attach", "exec", "proxy", "log", "portforward"),
	)

	kubeVersion := version.Get()
	genericConfig.Version = &kubeVersion

	storageFactoryConfig := kubeapiserver.NewStorageFactoryConfig(genericcontrolplanescheme.Scheme, genericcontrolplanescheme.Codecs)
	storageFactoryConfig.APIResourceConfig = genericConfig.MergedResourceConfig
	completedStorageFactoryConfig, err := storageFactoryConfig.Complete(s.Etcd)
	if err != nil {
		lastErr = err
		return
	}
	storageFactory, lastErr = completedStorageFactoryConfig.New()
	if lastErr != nil {
		return
	}
	if genericConfig.EgressSelector != nil {
		storageFactory.StorageConfig.Transport.EgressLookup = genericConfig.EgressSelector.Lookup
	}
	if utilfeature.DefaultFeatureGate.Enabled(genericfeatures.APIServerTracing) && genericConfig.TracerProvider != nil {
		storageFactory.StorageConfig.Transport.TracerProvider = genericConfig.TracerProvider
	}
	if lastErr = s.Etcd.ApplyWithStorageFactoryTo(storageFactory, genericConfig); lastErr != nil {
		return
	}

	// Use protobufs for self-communication.
	// Since not every generic apiserver has to support protobufs, we
	// cannot default to it in generic apiserver and need to explicitly
	// set it in kube-apiserver.
	genericConfig.LoopbackClientConfig.ContentConfig.ContentType = "application/vnd.kubernetes.protobuf"
	// Disable compression for self-communication, since we are going to be
	// on a fast local network
	genericConfig.LoopbackClientConfig.DisableCompression = true

	kubeClientConfig := genericConfig.LoopbackClientConfig

	clientutils.EnableMultiCluster(genericConfig.LoopbackClientConfig, genericConfig, true, "namespaces", "apiservices", "customresourcedefinitions")

	clientgoExternalClient, err := clientgoclientset.NewForConfig(kubeClientConfig)
	if err != nil {
		lastErr = fmt.Errorf("failed to create real external clientset: %v", err)
		return
	}
	versionedInformers := clientgoinformers.NewSharedInformerFactory(clientgoExternalClient, 10*time.Minute)

	// Authentication.ApplyTo requires already applied OpenAPIConfig and EgressSelector if present
	if lastErr = AuthenticationApplyTo(s.Authentication, &genericConfig.Authentication, genericConfig.SecureServing, genericConfig.EgressSelector, genericConfig.OpenAPIConfig); lastErr != nil {
		return
	}

	genericConfig.Authorization.Authorizer, genericConfig.RuleResolver, err = BuildAuthorizer(s, versionedInformers)
	if err != nil {
		lastErr = fmt.Errorf("invalid authorization config: %v", err)
		return
	}

	return
}

// BuildAuthorizer constructs the authorizer
func BuildAuthorizer(s *options.ServerRunOptions, versionedInformers clientgoinformers.SharedInformerFactory) (authorizer.Authorizer, authorizer.RuleResolver, error) {
	var (
		authorizers   []authorizer.Authorizer
		ruleResolvers []authorizer.RuleResolver
	)

	rbacAuthorizer := rbac.New(
		&rbac.RoleGetter{Lister: versionedInformers.Rbac().V1().Roles().Lister()},
		&rbac.RoleBindingLister{Lister: versionedInformers.Rbac().V1().RoleBindings().Lister()},
		&rbac.ClusterRoleGetter{Lister: versionedInformers.Rbac().V1().ClusterRoles().Lister()},
		&rbac.ClusterRoleBindingLister{Lister: versionedInformers.Rbac().V1().ClusterRoleBindings().Lister()},
	)
	authorizers = append(authorizers, rbacAuthorizer)
	ruleResolvers = append(ruleResolvers, rbacAuthorizer)

	return union.New(authorizers...), union.NewRuleResolvers(ruleResolvers...), nil
}

// BuildPriorityAndFairness constructs the guts of the API Priority and Fairness filter
func BuildPriorityAndFairness(s *options.ServerRunOptions, extclient clientgoclientset.Interface, versionedInformer clientgoinformers.SharedInformerFactory) (utilflowcontrol.Interface, error) {
	if s.GenericServerRunOptions.MaxRequestsInFlight+s.GenericServerRunOptions.MaxMutatingRequestsInFlight <= 0 {
		return nil, fmt.Errorf("invalid configuration: MaxRequestsInFlight=%d and MaxMutatingRequestsInFlight=%d; they must add up to something positive", s.GenericServerRunOptions.MaxRequestsInFlight, s.GenericServerRunOptions.MaxMutatingRequestsInFlight)
	}
	return utilflowcontrol.New(
		versionedInformer,
		extclient.FlowcontrolV1beta2(),
		s.GenericServerRunOptions.MaxRequestsInFlight+s.GenericServerRunOptions.MaxMutatingRequestsInFlight,
		s.GenericServerRunOptions.RequestTimeout/4,
	), nil
}

// completedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type completedServerRunOptions struct {
	*options.ServerRunOptions
}

// Complete set default ServerRunOptions.
// Should be called after kube-apiserver flags parsed.
func Complete(s *options.ServerRunOptions) (completedServerRunOptions, error) {
	var options completedServerRunOptions
	// set defaults
	if err := s.GenericServerRunOptions.DefaultAdvertiseAddress(s.SecureServing.SecureServingOptions); err != nil {
		return options, err
	}

	if err := s.SecureServing.MaybeDefaultWithSelfSignedCerts(s.GenericServerRunOptions.AdvertiseAddress.String(), nil, nil); err != nil {
		return options, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	if len(s.GenericServerRunOptions.ExternalHost) == 0 {
		if len(s.GenericServerRunOptions.AdvertiseAddress) > 0 {
			s.GenericServerRunOptions.ExternalHost = s.GenericServerRunOptions.AdvertiseAddress.String()
		} else {
			if hostname, err := os.Hostname(); err == nil {
				s.GenericServerRunOptions.ExternalHost = hostname
			} else {
				return options, fmt.Errorf("error finding host name: %v", err)
			}
		}
		klog.Infof("external host was not specified, using %v", s.GenericServerRunOptions.ExternalHost)
	}

	for key, value := range s.APIEnablement.RuntimeConfig {
		if key == "v1" || strings.HasPrefix(key, "v1/") ||
			key == "api/v1" || strings.HasPrefix(key, "api/v1/") {
			delete(s.APIEnablement.RuntimeConfig, key)
			s.APIEnablement.RuntimeConfig["/v1"] = value
		}
		if key == "api/legacy" {
			delete(s.APIEnablement.RuntimeConfig, key)
		}
	}

	options.ServerRunOptions = s
	return options, nil
}
