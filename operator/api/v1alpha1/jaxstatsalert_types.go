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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlertRule defines a single alerting rule
type AlertRule struct {
	// Name is the name of the alert rule
	Name string `json:"name"`

	// Metric is the metric to monitor
	Metric string `json:"metric"`

	// Operator is the comparison operator (e.g., ">", "<", "==")
	Operator string `json:"operator"`

	// Threshold is the value to compare against
	Threshold float64 `json:"threshold"`

	// Duration is how long the condition must be true before alerting
	Duration string `json:"duration,omitempty"`

	// Severity is the alert severity (e.g., "critical", "warning", "info")
	Severity string `json:"severity"`
}

// NotificationConfig defines how alerts should be notified
type NotificationConfig struct {
	// Type is the notification type (e.g., "email", "slack", "webhook")
	Type string `json:"type"`

	// Endpoint is the notification endpoint
	Endpoint string `json:"endpoint,omitempty"`

	// SecretName is the name of the secret containing notification credentials
	SecretName string `json:"secretName,omitempty"`

	// Template is the notification template to use
	Template string `json:"template,omitempty"`
}

// JAXStatsAlertSpec defines the desired state of JAXStatsAlert.
type JAXStatsAlertSpec struct {
	// Rules defines the alerting rules
	Rules []AlertRule `json:"rules"`

	// Notifications defines how alerts should be notified
	Notifications []NotificationConfig `json:"notifications"`

	// CollectorRef references the JAXStatsCollector to monitor
	CollectorRef string `json:"collectorRef"`

	// Enabled determines whether the alert is enabled
	Enabled bool `json:"enabled"`

	// CooldownPeriod is the time to wait before sending another alert
	CooldownPeriod string `json:"cooldownPeriod,omitempty"`
}

// JAXStatsAlertStatus defines the observed state of JAXStatsAlert.
type JAXStatsAlertStatus struct {
	// LastAlertTime is the timestamp of the last alert
	LastAlertTime *metav1.Time `json:"lastAlertTime,omitempty"`

	// AlertStatus indicates the current status of the alert
	AlertStatus string `json:"alertStatus,omitempty"`

	// ActiveAlerts tracks currently active alerts
	ActiveAlerts []string `json:"activeAlerts,omitempty"`

	// ErrorMessage contains any error message if alerting failed
	ErrorMessage string `json:"errorMessage,omitempty"`

	// Conditions represent the latest available observations of the alert's current state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// JAXStatsAlert is the Schema for the jaxstatsalerts API
type JAXStatsAlert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JAXStatsAlertSpec   `json:"spec,omitempty"`
	Status JAXStatsAlertStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JAXStatsAlertList contains a list of JAXStatsAlert
type JAXStatsAlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JAXStatsAlert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JAXStatsAlert{}, &JAXStatsAlertList{})
}
