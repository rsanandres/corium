/*
Copyright 2025.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JAXStatsConfigSpec defines the desired state of JAXStatsConfig.
type JAXStatsConfigSpec struct {
	// Enabled determines whether JAXStats collection is enabled
	Enabled bool `json:"enabled"`

	// CollectionInterval specifies how often stats should be collected (in seconds)
	// +kubebuilder:validation:Minimum=10
	// +kubebuilder:validation:Maximum=3600
	// +kubebuilder:default=60
	CollectionInterval int32 `json:"collectionInterval,omitempty"`

	// Metrics defines which metrics should be collected
	// +kubebuilder:validation:MinItems=1
	Metrics []string `json:"metrics,omitempty"`

	// StorageConfig defines where the stats should be stored
	StorageConfig StorageConfig `json:"storageConfig,omitempty"`
}

// StorageConfig defines the storage configuration for JAXStats
type StorageConfig struct {
	// Type specifies the type of storage (e.g., "prometheus", "configmap")
	// +kubebuilder:validation:Enum=prometheus;configmap;elasticsearch
	Type string `json:"type"`

	// Endpoint is the endpoint URL for the storage backend
	Endpoint string `json:"endpoint,omitempty"`

	// CredentialsSecret is the name of the secret containing storage credentials
	CredentialsSecret string `json:"credentialsSecret,omitempty"`
}

// JAXStatsConfigStatus defines the observed state of JAXStatsConfig.
type JAXStatsConfigStatus struct {
	// LastCollectionTime is the timestamp of the last successful stats collection
	LastCollectionTime *metav1.Time `json:"lastCollectionTime,omitempty"`

	// CollectionStatus indicates the current status of stats collection
	CollectionStatus string `json:"collectionStatus,omitempty"`

	// ErrorMessage contains any error message if collection failed
	ErrorMessage string `json:"errorMessage,omitempty"`

	// DependentCollectors tracks the number of Collectors referencing this Config
	DependentCollectors int32 `json:"dependentCollectors,omitempty"`

	// Conditions represent the latest available observations of the config's current state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=jsc
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.spec.enabled`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.collectionStatus`
// +kubebuilder:printcolumn:name="Interval",type=integer,JSONPath=`.spec.collectionInterval`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// JAXStatsConfig is the Schema for the jaxstatsconfigs API.
type JAXStatsConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JAXStatsConfigSpec   `json:"spec,omitempty"`
	Status JAXStatsConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JAXStatsConfigList contains a list of JAXStatsConfig.
type JAXStatsConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JAXStatsConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JAXStatsConfig{}, &JAXStatsConfigList{})
}
