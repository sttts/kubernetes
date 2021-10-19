/*
Copyright 2018 The Kubernetes Authors.

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

package kcp

import (
	"flag"
	"fmt"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/test/e2e/framework"
)

var (
	kcpMultiClusterKubeconfig = flag.String(fmt.Sprintf("%s-%s", "kcp-multi-cluster", clientcmd.RecommendedConfigPathFlag), "", "Path to kubeconfig containing embedded authinfo for kcp multi-cluster client.")
	kcpMultiClusterContext    = flag.String(fmt.Sprintf("%s-context", "kcp-multi-cluster"), "", "Context within multi-cluster kubeconfig to use.")

	kcpSecondaryKubeconfig = flag.String(fmt.Sprintf("%s-%s", "kcp-secondary", clientcmd.RecommendedConfigPathFlag), "", "Path to kubeconfig containing embedded authinfo for kcp secondary client.")
	kcpSecondaryContext = flag.String(fmt.Sprintf("%s-context", "kcp-secondary"), "", "Context within secondary kubeconfig to use.")

	kcpTertiaryKubeconfig = flag.String(fmt.Sprintf("%s-%s", "kcp-tertiary", clientcmd.RecommendedConfigPathFlag), "", "Path to kubeconfig containing embedded authinfo for kcp tertiary client.")
	kcpTertiaryContext = flag.String(fmt.Sprintf("%s-context", "kcp-tertiary"), "", "Context within tertiary kubeconfig to use.")

	kcpClusterlessKubeconfig = flag.String(fmt.Sprintf("%s-%s", "kcp-clusterless", clientcmd.RecommendedConfigPathFlag), "", "Path to kubeconfig containing embedded authinfo for kcp clusterless client.")
	kcpClusterlessContext = flag.String(fmt.Sprintf("%s-context", "kcp-clusterless"), "", "Context within clusterless kubeconfig to use.")
)

func init() {
	framework.RegisterProvider("kcp", newProvider)
}

func newProvider() (framework.ProviderInterface, error) {
	// Actual initialization happens when the e2e framework gets constructed.
	return &Provider{}, nil
}

// Provider is a structure to handle Kubemark cluster for e2e testing
type Provider struct {
	framework.NullProvider
}

func clientSetFor(kubeconfig, context string) clientset.Interface {
	rawConfig, err := clientcmd.LoadFromFile(kubeconfig)
	framework.ExpectNoError(err)
	config, err := clientcmd.NewNonInteractiveClientConfig(*rawConfig, context, nil, nil).ClientConfig()
	framework.ExpectNoError(err)
	client, err := clientset.NewForConfig(config)
	framework.ExpectNoError(err)
	return client
}

// FrameworkBeforeEach prepares clients, configurations etc. for e2e testing
func (p *Provider) FrameworkBeforeEach(f *framework.Framework) {
	if *kcpMultiClusterKubeconfig != "" && *kcpMultiClusterContext != "" && f.MultiClusterClientSet == nil {
		f.MultiClusterClientSet = clientSetFor(*kcpMultiClusterKubeconfig, *kcpMultiClusterContext)
	}
	if *kcpSecondaryKubeconfig != "" && *kcpSecondaryContext != "" && f.SecondaryClientSet == nil {
		f.SecondaryClientSet = clientSetFor(*kcpSecondaryKubeconfig, *kcpSecondaryContext)
	}
	if *kcpTertiaryKubeconfig != "" && *kcpTertiaryContext != "" && f.TertiaryClientSet == nil {
		f.TertiaryClientSet = clientSetFor(*kcpTertiaryKubeconfig, *kcpTertiaryContext)
	}
	if *kcpClusterlessKubeconfig != "" && *kcpClusterlessContext != "" && f.ClusterlessClientSet == nil {
		rawConfig, err := clientcmd.LoadFromFile(*kcpClusterlessKubeconfig)
		framework.ExpectNoError(err)
		config, err := clientcmd.NewNonInteractiveClientConfig(*rawConfig, *kcpClusterlessContext, nil, nil).ClientConfig()
		framework.ExpectNoError(err)
		client, err := clientset.NewClusterForConfig(config)
		framework.ExpectNoError(err)
		f.ClusterlessClientSet = client
	}
}
