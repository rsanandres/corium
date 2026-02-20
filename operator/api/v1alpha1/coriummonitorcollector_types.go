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

// CoriumMonitorCollectorSpec defines the desired state of CoriumMonitorCollector.
type CoriumMonitorCollectorSpec struct {
	// TargetNamespace specifies which namespace to collect stats from
	// +kubebuilder:validation:MinLength=1
	TargetNamespace string `json:"targetNamespace"`

	// Selector specifies which resources to collect stats from
	Selector metav1.LabelSelector `json:"selector"`

	// Metrics specifies which metrics to collect
	// +kubebuilder:validation:MinItems=1
	Metrics []string `json:"metrics"`

	// ConfigRef references the CoriumMonitorConfig to use
	// +kubebuilder:validation:MinLength=1
	ConfigRef string `json:"configRef"`

	// ResourceTypes specifies which types of resources to collect stats from
	ResourceTypes []string `json:"resourceTypes,omitempty"`

	// CollectionSchedule specifies when to collect stats (cron format)
	CollectionSchedule string `json:"collectionSchedule,omitempty"`
}

// CoriumMonitorCollectorStatus defines the observed state of CoriumMonitorCollector.
type CoriumMonitorCollectorStatus struct {
	// LastCollectionTime is the timestamp of the last successful collection
	LastCollectionTime *metav1.Time `json:"lastCollectionTime,omitempty"`

	// CollectionStatus indicates the current status of collection
	CollectionStatus string `json:"collectionStatus,omitempty"`

	// ErrorMessage contains any error message if collection failed
	ErrorMessage string `json:"errorMessage,omitempty"`

	// CollectedResources tracks the number of resources being monitored
	CollectedResources int32 `json:"collectedResources,omitempty"`

	// DiscoveredPods tracks the number of pods discovered by the selector
	DiscoveredPods int32 `json:"discoveredPods,omitempty"`

	// MetricsConfigMap is the name of the ConfigMap storing collected metrics
	MetricsConfigMap string `json:"metricsConfigMap,omitempty"`

	// Conditions represent the latest available observations of the collector's current state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=cmcol
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.collectionStatus`
// +kubebuilder:printcolumn:name="Pods",type=integer,JSONPath=`.status.discoveredPods`
// +kubebuilder:printcolumn:name="Config",type=string,JSONPath=`.spec.configRef`
// +kubebuilder:printcolumn:name="Namespace",type=string,JSONPath=`.spec.targetNamespace`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// CoriumMonitorCollector is the Schema for the coriummonitorcollectors API.
type CoriumMonitorCollector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoriumMonitorCollectorSpec   `json:"spec,omitempty"`
	Status CoriumMonitorCollectorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CoriumMonitorCollectorList contains a list of CoriumMonitorCollector.
type CoriumMonitorCollectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoriumMonitorCollector `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CoriumMonitorCollector{}, &CoriumMonitorCollectorList{})
}
