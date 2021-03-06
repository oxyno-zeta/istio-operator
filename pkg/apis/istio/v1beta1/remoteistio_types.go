/*
Copyright 2019 Banzai Cloud.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SignCert struct {
	CA    []byte
	Root  []byte
	Key   []byte
	Chain []byte
}

type IstioService struct {
	Name          string   `json:"name"`
	LabelSelector string   `json:"labelSelector"`
	IPs           []string `json:"podIPs,omitempty"`
}

func (spec RemoteIstioSpec) SetSignCert(signCert SignCert) RemoteIstioSpec {
	spec.signCert = signCert
	return spec
}

func (spec RemoteIstioSpec) GetSignCert() SignCert {
	return spec.signCert
}

// RemoteIstioSpec defines the desired state of RemoteIstio
type RemoteIstioSpec struct {
	IncludeIPRanges string         `json:"includeIPRanges,omitempty"`
	ExcludeIPRanges string         `json:"excludeIPRanges,omitempty"`
	EnabledServices []IstioService `json:"enabledServices"`
	// List of namespaces to label with sidecar auto injection enabled
	AutoInjectionNamespaces []string `json:"autoInjectionNamespaces,omitempty"`
	// ControlPlaneSecurityEnabled control plane services are communicating through mTLS
	ControlPlaneSecurityEnabled bool `json:"controlPlaneSecurityEnabled,omitempty"`
	// Citadel configuration options
	Citadel CitadelConfiguration `json:"citadel,omitempty"`
	// SidecarInjector configuration options
	SidecarInjector SidecarInjectorConfiguration `json:"sidecarInjector,omitempty"`

	signCert SignCert
}

// RemoteIstioStatus defines the observed state of RemoteIstio
type RemoteIstioStatus struct {
	Status       ConfigState
	ErrorMessage string
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RemoteIstio is the Schema for the remoteistios API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type RemoteIstio struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RemoteIstioSpec   `json:"spec,omitempty"`
	Status RemoteIstioStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RemoteIstioList contains a list of RemoteIstio
type RemoteIstioList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RemoteIstio `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RemoteIstio{}, &RemoteIstioList{})
}

func SetRemoteIstioDefaults(remoteconfig *RemoteIstio) {
	if remoteconfig.Spec.IncludeIPRanges == "" {
		remoteconfig.Spec.IncludeIPRanges = defaultIncludeIPRanges
	}
	// Citadel config
	if remoteconfig.Spec.Citadel.ReplicaCount == 0 {
		remoteconfig.Spec.Citadel.ReplicaCount = defaultReplicaCount
	}
	// SidecarInjector config
	if remoteconfig.Spec.SidecarInjector.ReplicaCount == 0 {
		remoteconfig.Spec.SidecarInjector.ReplicaCount = defaultReplicaCount
	}
}
