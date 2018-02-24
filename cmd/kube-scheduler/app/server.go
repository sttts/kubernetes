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

// Package app implements a Server object for running the scheduler.
package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	goruntime "runtime"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/apiserver/pkg/server/routes"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	storageinformers "k8s.io/client-go/informers/storage/v1"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection"
	schedulerserverconfig "k8s.io/kubernetes/cmd/kube-scheduler/app/config"
	"k8s.io/kubernetes/cmd/kube-scheduler/app/options"
	"k8s.io/kubernetes/pkg/apis/componentconfig"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/scheduler"
	"k8s.io/kubernetes/pkg/scheduler/algorithmprovider"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	latestschedulerapi "k8s.io/kubernetes/pkg/scheduler/api/latest"
	"k8s.io/kubernetes/pkg/scheduler/factory"
	"k8s.io/kubernetes/pkg/util/configz"
	utilflag "k8s.io/kubernetes/pkg/util/flag"
	"k8s.io/kubernetes/pkg/version"
	"k8s.io/kubernetes/pkg/version/verflag"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

// NewSchedulerCommand creates a *cobra.Command object with default parameters
func NewSchedulerCommand() *cobra.Command {
	opts, err := options.NewOptions()
	if err != nil {
		glog.Fatalf("unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use: "kube-scheduler",
		Long: `The Kubernetes scheduler is a policy-rich, topology-aware,
workload-specific function that significantly impacts availability, performance,
and capacity. The scheduler needs to take into account individual and collective
resource requirements, quality of service requirements, hardware/software/policy
constraints, affinity and anti-affinity specifications, data locality, inter-workload
interference, deadlines, and so on. Workload-specific requirements will be exposed
through the API as necessary.`,
		Run: func(cmd *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()
			utilflag.PrintFlags(cmd.Flags())

			if len(args) != 0 {
				fmt.Fprint(os.Stderr, "arguments are not supported\n")
			}

			c, err := opts.Config()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			stopCh := make(chan struct{})
			if err := Run(c.Complete(), stopCh); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	opts.AddFlags(cmd.Flags())
	cmd.MarkFlagFilename("config", "yaml", "yml", "json")

	return cmd
}

func Run(c schedulerserverconfig.CompletedConfig, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	glog.Infof("Version: %+v", version.Get())

	// Apply algorithms based on feature gates.
	// TODO: make configurable?
	algorithmprovider.ApplyFeatureGates()

	// Configz registration.
	if cz, err := configz.New("componentconfig"); err == nil {
		cz.Set(c.ComponentConfig)
	} else {
		return fmt.Errorf("unable to register configz: %s", err)
	}

	// Build a scheduler config from the provided algorithm source.
	schedulerConfig, err := NewSchedulerConfig(c)
	if err != nil {
		return err
	}

	// Create the scheduler.
	sched := scheduler.NewFromConfig(schedulerConfig)

	// Prepare the event broadcaster.
	if c.Broadcaster != nil && c.EventClient != nil {
		c.Broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: c.EventClient.Events("")})
	}

	// Start up the healthz server.
	if c.InsecureServing != nil {
		handler := buildHandlerChain(newHealthzHandler(&c.ComponentConfig, c.InsecureMetricsServing != nil))
		// TODO: fail early as all other Kubernetes binaries
		go wait.Until(func() {
			if err := c.InsecureServing.Serve(handler, 0, stopCh); err != nil {
				utilruntime.HandleError(fmt.Errorf("failed to start healthz server: %v", err))
			}
		}, 5*time.Second, stopCh)
	}
	if c.InsecureServing != nil {
		handler := buildHandlerChain(newMetricsHandler(&c.ComponentConfig))
		// TODO: fail early as all other Kubernetes binaries
		go wait.Until(func() {
			if err := c.InsecureServing.Serve(handler, 0, stopCh); err != nil {
				utilruntime.HandleError(fmt.Errorf("failed to start metrics server: %v", err))
			}
		}, 5*time.Second, stopCh)
	}

	// Start all informers.
	go c.PodInformer.Informer().Run(stopCh)
	c.InformerFactory.Start(stopCh)

	// Wait for all caches to sync before scheduling.
	c.InformerFactory.WaitForCacheSync(stopCh)
	controller.WaitForCacheSync("scheduler", stopCh, c.PodInformer.Informer().HasSynced)

	// Prepare a reusable run function.
	run := func(stopCh <-chan struct{}) {
		sched.Run()
		<-stopCh
	}

	// If leader election is enabled, run via LeaderElector until done and exit.
	if c.LeaderElection != nil {
		c.LeaderElection.Callbacks = leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				utilruntime.HandleError(fmt.Errorf("lost master"))
			},
		}
		leaderElector, err := leaderelection.NewLeaderElector(*c.LeaderElection)
		if err != nil {
			return fmt.Errorf("couldn't create leader elector: %v", err)
		}

		leaderElector.Run()

		return fmt.Errorf("lost lease")
	}

	// Leader election is disabled, so run inline until done.
	run(stopCh)
	return fmt.Errorf("finished without leader elect")
}

// buildHandlerChain wraps the given handler with the standard filters.
func buildHandlerChain(handler http.Handler) http.Handler {
	requestContextMapper := apirequest.NewRequestContextMapper()
	requestInfoResolver := &apirequest.RequestInfoFactory{}

	handler = genericapifilters.WithRequestInfo(handler, requestInfoResolver, requestContextMapper)
	handler = apirequest.WithRequestContext(handler, requestContextMapper)
	handler = genericfilters.WithPanicRecovery(handler)

	return handler
}

// newMetricsHandler builds a metrics server from the config.
func newMetricsHandler(config *componentconfig.KubeSchedulerConfiguration) http.Handler {
	mux := mux.NewPathRecorderMux("kube-scheduler")
	configz.InstallHandler(mux)
	mux.Handle("/metrics", prometheus.Handler())
	if config.EnableProfiling {
		routes.Profiling{}.Install(mux)
		if config.EnableContentionProfiling {
			goruntime.SetBlockProfileRate(1)
		}
	}
	return mux
}

// newHealthzServer creates a healthz server from the config, and will also
// embed the metrics handler if the healthz and metrics address configurations
// are the same.
func newHealthzHandler(config *componentconfig.KubeSchedulerConfiguration, separateMetrics bool) http.Handler {
	mux := mux.NewPathRecorderMux("kube-scheduler")
	healthz.InstallHandler(mux)
	if !separateMetrics {
		configz.InstallHandler(mux)
		mux.Handle("/metrics", prometheus.Handler())
	}
	if config.EnableProfiling {
		routes.Profiling{}.Install(mux)
		if config.EnableContentionProfiling {
			goruntime.SetBlockProfileRate(1)
		}
	}
	return mux
}

// NewSchedulerConfig creates the scheduler configuration. This is exposed for use by tests.
func NewSchedulerConfig(s schedulerserverconfig.CompletedConfig) (*scheduler.Config, error) {
	var storageClassInformer storageinformers.StorageClassInformer
	if utilfeature.DefaultFeatureGate.Enabled(features.VolumeScheduling) {
		storageClassInformer = s.InformerFactory.Storage().V1().StorageClasses()
	}

	// Set up the configurator which can create schedulers from configs.
	configurator := factory.NewConfigFactory(
		s.SchedulerName,
		s.Client,
		s.InformerFactory.Core().V1().Nodes(),
		s.PodInformer,
		s.InformerFactory.Core().V1().PersistentVolumes(),
		s.InformerFactory.Core().V1().PersistentVolumeClaims(),
		s.InformerFactory.Core().V1().ReplicationControllers(),
		s.InformerFactory.Extensions().V1beta1().ReplicaSets(),
		s.InformerFactory.Apps().V1beta1().StatefulSets(),
		s.InformerFactory.Core().V1().Services(),
		s.InformerFactory.Policy().V1beta1().PodDisruptionBudgets(),
		storageClassInformer,
		s.HardPodAffinitySymmetricWeight,
		utilfeature.DefaultFeatureGate.Enabled(features.EnableEquivalenceClassCache),
	)

	source := s.AlgorithmSource
	var config *scheduler.Config
	switch {
	case source.Provider != nil:
		// Create the config from a named algorithm provider.
		sc, err := configurator.CreateFromProvider(*source.Provider)
		if err != nil {
			return nil, fmt.Errorf("couldn't create scheduler using provider %q: %v", *source.Provider, err)
		}
		config = sc
	case source.Policy != nil:
		// Create the config from a user specified policy source.
		policy := &schedulerapi.Policy{}
		switch {
		case source.Policy.File != nil:
			// Use a policy serialized in a file.
			policyFile := source.Policy.File.Path
			_, err := os.Stat(policyFile)
			if err != nil {
				return nil, fmt.Errorf("missing policy config file %s", policyFile)
			}
			data, err := ioutil.ReadFile(policyFile)
			if err != nil {
				return nil, fmt.Errorf("couldn't read policy config: %v", err)
			}
			err = runtime.DecodeInto(latestschedulerapi.Codec, []byte(data), policy)
			if err != nil {
				return nil, fmt.Errorf("invalid policy: %v", err)
			}
		case source.Policy.ConfigMap != nil:
			// Use a policy serialized in a config map value.
			policyRef := source.Policy.ConfigMap
			policyConfigMap, err := s.Client.CoreV1().ConfigMaps(policyRef.Namespace).Get(policyRef.Name, metav1.GetOptions{})
			if err != nil {
				return nil, fmt.Errorf("couldn't get policy config map %s/%s: %v", policyRef.Namespace, policyRef.Name, err)
			}
			data, found := policyConfigMap.Data[componentconfig.SchedulerPolicyConfigMapKey]
			if !found {
				return nil, fmt.Errorf("missing policy config map value at key %q", componentconfig.SchedulerPolicyConfigMapKey)
			}
			err = runtime.DecodeInto(latestschedulerapi.Codec, []byte(data), policy)
			if err != nil {
				return nil, fmt.Errorf("invalid policy: %v", err)
			}
		}
		sc, err := configurator.CreateFromConfig(*policy)
		if err != nil {
			return nil, fmt.Errorf("couldn't create scheduler from policy: %v", err)
		}
		config = sc
	default:
		return nil, fmt.Errorf("unsupported algorithm source: %v", source)
	}
	// Additional tweaks to the config produced by the configurator.
	config.Recorder = s.Recorder
	return config, nil
}
