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

package app

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	goruntime "runtime"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	cloudcontrollers "k8s.io/kubernetes/pkg/controller/cloud"
	routecontroller "k8s.io/kubernetes/pkg/controller/route"
	servicecontroller "k8s.io/kubernetes/pkg/controller/service"
	"k8s.io/kubernetes/pkg/util/configz"
	"k8s.io/kubernetes/pkg/version/verflag"
)

const (
	// ControllerStartJitter is the jitter value used when starting controller managers.
	ControllerStartJitter = 1.0
)

// NewCloudControllerManagerCommand creates a *cobra.Command object with default parameters
func NewCloudControllerManagerCommand() *cobra.Command {
	s := options.NewCloudControllerManagerOptions()
	cmd := &cobra.Command{
		Use: "cloud-controller-manager",
		Long: `The Cloud controller manager is a daemon that embeds
the cloud specific control loops shipped with Kubernetes.`,
		Run: func(cmd *cobra.Command, args []string) {
			verflag.PrintAndExitIfRequested()

			if err := Run(s); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

		},
	}
	s.AddFlags(cmd.Flags())

	return cmd
}

// resyncPeriod computes the time interval a shared informer waits before resyncing with the api server
func resyncPeriod(s *options.CloudControllerManagerOptions) func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(s.Generic.ComponentConfig.MinResyncPeriod.Nanoseconds()) * factor)
	}
}

// Run runs the ExternalCMServer.  This should never exit.
func Run(s *CloudControllerManagerConfig) error {
	if s.Generic.ComponentConfig.CloudProvider == "" {
		glog.Fatalf("--cloud-provider cannot be empty")
	}

	cloud, err := cloudprovider.InitCloudProvider(s.Generic.ComponentConfig.CloudProvider, s.Generic.ComponentConfig.CloudConfigFile)
	if err != nil {
		glog.Fatalf("Cloud provider could not be initialized: %v", err)
	}

	if cloud == nil {
		glog.Fatalf("cloud provider is nil")
	}

	if cloud.HasClusterID() == false {
		if s.Generic.ComponentConfig.AllowUntaggedCloud == true {
			glog.Warning("detected a cluster without a ClusterID.  A ClusterID will be required in the future.  Please tag your cluster to avoid any future issues")
		} else {
			glog.Fatalf("no ClusterID found.  A ClusterID is required for the cloud provider to function properly.  This check can be bypassed by setting the allow-untagged-cloud option")
		}
	}

	if c, err := configz.New("componentconfig"); err == nil {
		c.Set(s.ControllerManagerOptions)
	} else {
		glog.Errorf("unable to register configz: %s", err)
	}
	kubeconfig, err := clientcmd.BuildConfigFromFlags(s.Generic.Master, s.Generic.Kubeconfig)
	if err != nil {
		return err
	}

	// Set the ContentType of the requests from kube client
	kubeconfig.ContentConfig.ContentType = s.Generic.ComponentConfig.ContentType
	// Override kubeconfig qps/burst settings from flags
	kubeconfig.QPS = s.Generic.ComponentConfig.KubeAPIQPS
	kubeconfig.Burst = int(s.Generic.ComponentConfig.KubeAPIBurst)
	kubeClient, err := kubernetes.NewForConfig(restclient.AddUserAgent(kubeconfig, "cloud-controller-manager"))
	if err != nil {
		glog.Fatalf("Invalid API configuration: %v", err)
	}
	leaderElectionClient := kubernetes.NewForConfigOrDie(restclient.AddUserAgent(kubeconfig, "leader-election"))

	// Start the external controller manager server
	go startHTTP(s)

	recorder := createRecorder(kubeClient)

	run := func(stop <-chan struct{}) {
		rootClientBuilder := controller.SimpleControllerClientBuilder{
			ClientConfig: kubeconfig,
		}
		var clientBuilder controller.ControllerClientBuilder
		if s.Generic.ComponentConfig.UseServiceAccountCredentials {
			clientBuilder = controller.SAControllerClientBuilder{
				ClientConfig:         restclient.AnonymousClientConfig(kubeconfig),
				CoreClient:           kubeClient.CoreV1(),
				AuthenticationClient: kubeClient.AuthenticationV1(),
				Namespace:            "kube-system",
			}
		} else {
			clientBuilder = rootClientBuilder
		}

		if err := StartControllers(s, kubeconfig, rootClientBuilder, clientBuilder, stop, recorder, cloud); err != nil {
			glog.Fatalf("error running controllers: %v", err)
		}
	}

	if !s.Generic.ComponentConfig.LeaderElection.LeaderElect {
		run(nil)
		panic("unreachable")
	}

	// Identity used to distinguish between multiple cloud controller manager instances
	id, err := os.Hostname()
	if err != nil {
		return err
	}
	// add a uniquifier so that two processes on the same host don't accidentally both become active
	id = id + "_" + string(uuid.NewUUID())

	// Lock required for leader election
	rl, err := resourcelock.New(s.Generic.ComponentConfig.LeaderElection.ResourceLock,
		"kube-system",
		"cloud-controller-manager",
		leaderElectionClient.CoreV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: recorder,
		})
	if err != nil {
		glog.Fatalf("error creating lock: %v", err)
	}

	// Try and become the leader and start cloud controller manager loops
	leaderelection.RunOrDie(leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: s.Generic.ComponentConfig.LeaderElection.LeaseDuration.Duration,
		RenewDeadline: s.Generic.ComponentConfig.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:   s.Generic.ComponentConfig.LeaderElection.RetryPeriod.Duration,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				glog.Fatalf("leaderelection lost")
			},
		},
	})
	panic("unreachable")
}

// StartControllers starts the cloud specific controller loops.
func StartControllers(s *options.CloudControllerManagerOptions, kubeconfig *restclient.Config, rootClientBuilder, clientBuilder controller.ControllerClientBuilder, stop <-chan struct{}, recorder record.EventRecorder, cloud cloudprovider.Interface) error {
	// Function to build the kube client object
	client := func(serviceAccountName string) kubernetes.Interface {
		return clientBuilder.ClientOrDie(serviceAccountName)
	}

	if cloud != nil {
		// Initialize the cloud provider with a reference to the clientBuilder
		cloud.Initialize(clientBuilder)
	}

	versionedClient := rootClientBuilder.ClientOrDie("shared-informers")
	sharedInformers := informers.NewSharedInformerFactory(versionedClient, resyncPeriod(s)())

	// Start the CloudNodeController
	nodeController := cloudcontrollers.NewCloudNodeController(
		sharedInformers.Core().V1().Nodes(),
		client("cloud-node-controller"), cloud,
		s.Generic.ComponentConfig.NodeMonitorPeriod.Duration,
		s.NodeStatusUpdateFrequency.Duration)

	nodeController.Run()
	time.Sleep(wait.Jitter(s.Generic.ComponentConfig.ControllerStartInterval.Duration, ControllerStartJitter))

	// Start the PersistentVolumeLabelController
	pvlController := cloudcontrollers.NewPersistentVolumeLabelController(client("pvl-controller"), cloud)
	threads := 5
	go pvlController.Run(threads, stop)
	time.Sleep(wait.Jitter(s.Generic.ComponentConfig.ControllerStartInterval.Duration, ControllerStartJitter))

	// Start the service controller
	serviceController, err := servicecontroller.New(
		cloud,
		client("service-controller"),
		sharedInformers.Core().V1().Services(),
		sharedInformers.Core().V1().Nodes(),
		s.Generic.ComponentConfig.ClusterName,
	)
	if err != nil {
		glog.Errorf("Failed to start service controller: %v", err)
	} else {
		go serviceController.Run(stop, int(s.Generic.ComponentConfig.ConcurrentServiceSyncs))
		time.Sleep(wait.Jitter(s.Generic.ComponentConfig.ControllerStartInterval.Duration, ControllerStartJitter))
	}

	// If CIDRs should be allocated for pods and set on the CloudProvider, then start the route controller
	if s.Generic.ComponentConfig.AllocateNodeCIDRs && s.Generic.ComponentConfig.ConfigureCloudRoutes {
		if routes, ok := cloud.Routes(); !ok {
			glog.Warning("configure-cloud-routes is set, but cloud provider does not support routes. Will not configure cloud provider routes.")
		} else {
			var clusterCIDR *net.IPNet
			if len(strings.TrimSpace(s.Generic.ComponentConfig.ClusterCIDR)) != 0 {
				_, clusterCIDR, err = net.ParseCIDR(s.Generic.ComponentConfig.ClusterCIDR)
				if err != nil {
					glog.Warningf("Unsuccessful parsing of cluster CIDR %v: %v", s.Generic.ComponentConfig.ClusterCIDR, err)
				}
			}

			routeController := routecontroller.New(routes, client("route-controller"), sharedInformers.Core().V1().Nodes(), s.Generic.ComponentConfig.ClusterName, clusterCIDR)
			go routeController.Run(stop, s.Generic.ComponentConfig.RouteReconciliationPeriod.Duration)
			time.Sleep(wait.Jitter(s.Generic.ComponentConfig.ControllerStartInterval.Duration, ControllerStartJitter))
		}
	} else {
		glog.Infof("Will not configure cloud provider routes for allocate-node-cidrs: %v, configure-cloud-routes: %v.", s.Generic.ComponentConfig.AllocateNodeCIDRs, s.Generic.ComponentConfig.ConfigureCloudRoutes)
	}

	// If apiserver is not running we should wait for some time and fail only then. This is particularly
	// important when we start apiserver and controller manager at the same time.
	err = wait.PollImmediate(time.Second, 10*time.Second, func() (bool, error) {
		if _, err = restclient.ServerAPIVersions(kubeconfig); err == nil {
			return true, nil
		}
		glog.Errorf("Failed to get api versions from server: %v", err)
		return false, nil
	})
	if err != nil {
		glog.Fatalf("Failed to get api versions from server: %v", err)
	}

	sharedInformers.Start(stop)

	select {}
}

func startHTTP(s *options.CloudControllerManagerOptions) {
	mux := http.NewServeMux()
	healthz.InstallHandler(mux)
	if s.Generic.ComponentConfig.EnableProfiling {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		if s.Generic.ComponentConfig.EnableContentionProfiling {
			goruntime.SetBlockProfileRate(1)
		}
	}
	configz.InstallHandler(mux)
	mux.Handle("/metrics", prometheus.Handler())

	server := &http.Server{
		Addr:    net.JoinHostPort(s.Generic.ComponentConfig.Address, strconv.Itoa(int(s.Generic.ComponentConfig.Port))),
		Handler: mux,
	}
	glog.Fatal(server.ListenAndServe())
}

func createRecorder(kubeClient *kubernetes.Clientset) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.CoreV1().RESTClient()).Events("")})
	return eventBroadcaster.NewRecorder(legacyscheme.Scheme, v1.EventSource{Component: "cloud-controller-manager"})
}
