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

// +kcp-code-generator:skip

package serviceaccount

import (
	"context"

	kcpkubernetesclientset "github.com/kcp-dev/client-go/kubernetes"
	kcpcorev1listers "github.com/kcp-dev/client-go/listers/core/v1"
	"github.com/kcp-dev/logicalcluster/v3"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	v1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/kubernetes/pkg/serviceaccount"
)

func NewClusterGetterFromClient(client kcpkubernetesclientset.ClusterInterface, secretLister kcpcorev1listers.SecretClusterLister, serviceAccountLister kcpcorev1listers.ServiceAccountClusterLister /*podLister kcpcorev1listers.PodClusterLister*/) serviceaccount.ServiceAccountTokenClusterGetter {
	return &serviceAccountTokenClusterGetter{
		client:               client,
		secretLister:         secretLister,
		serviceAccountLister: serviceAccountLister,
	}
}

type serviceAccountTokenClusterGetter struct {
	client               kcpkubernetesclientset.ClusterInterface
	secretLister         kcpcorev1listers.SecretClusterLister
	serviceAccountLister kcpcorev1listers.ServiceAccountClusterLister
	podLister            kcpcorev1listers.PodClusterLister
}

func (s *serviceAccountTokenClusterGetter) Cluster(name logicalcluster.Name) serviceaccount.ServiceAccountTokenGetter {
	return NewGetterFromClient(
		s.client.Cluster(name),
		s.secretLister.Cluster(name),
		s.serviceAccountLister.Cluster(name),
		nil,
	)
}

// clientGetter implements ServiceAccountTokenGetter using a clientset.Interface
type clientGetter struct {
	client               clientset.Interface
	secretLister         v1listers.SecretLister
	serviceAccountLister v1listers.ServiceAccountLister
	podLister            v1listers.PodLister
}

// NewGetterFromClient returns a ServiceAccountTokenGetter that
// uses the specified client to retrieve service accounts and secrets.
// The client should NOT authenticate using a service account token
// the returned getter will be used to retrieve, or recursion will result.
func NewGetterFromClient(c clientset.Interface, secretLister v1listers.SecretLister, serviceAccountLister v1listers.ServiceAccountLister, podLister v1listers.PodLister) serviceaccount.ServiceAccountTokenGetter {
	return clientGetter{c, secretLister, serviceAccountLister, podLister}
}

func (c clientGetter) GetServiceAccount(namespace, name string) (*v1.ServiceAccount, error) {
	if serviceAccount, err := c.serviceAccountLister.ServiceAccounts(namespace).Get(name); err == nil {
		return serviceAccount, nil
	}
	return c.client.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c clientGetter) GetPod(namespace, name string) (*v1.Pod, error) {
	return c.client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c clientGetter) GetSecret(namespace, name string) (*v1.Secret, error) {
	if secret, err := c.secretLister.Secrets(namespace).Get(name); err == nil {
		return secret, nil
	}
	return c.client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
