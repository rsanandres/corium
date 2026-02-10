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

// AlertRule defines a single alerting rule
type AlertRule struct {
	// Name is the name of the alert rule
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// Metric is the metric to monitor
	// +kubebuilder:validation:Enum=restart_count;not_ready_count;container_count;pod_count
	Metric string `json:"metric"`

	// Operator is the comparison operator
	// +kubebuilder:validation:Enum=">";"<";">=";"<=";"=="
	Operator string `json:"operator"`

	// Threshold is the value to compare against
	Threshold string `json:"threshold"`

	// Duration is how long the condition must be true before alerting
	Duration string `json:"duration,omitempty"`

	// Severity is the alert severity
	// +kubebuilder:validation:Enum=critical;warning;info
	Severity string `json:"severity"`
}

// NotificationConfig defines how alerts should be notified
type NotificationConfig struct {
	// Type is the notification type (e.g., "event", "slack", "webhook")
	// +kubebuilder:validation:Enum=event;slack;webhook
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
	// +kubebuilder:validation:MinItems=1
	Rules []AlertRule `json:"rules"`

	// Notifications defines how alerts should be notified
	Notifications []NotificationConfig `json:"notifications,omitempty"`

	// CollectorRef references the JAXStatsCollector to monitor
	// +kubebuilder:validation:MinLength=1
	CollectorRef string `json:"collectorRef"`

	// Enabled determines whether the alert is enabled
	Enabled bool `json:"enabled"`

	// CooldownPeriod is the time to wait before sending another alert (e.g., "5m", "1h")
	// +kubebuilder:default="5m"
	CooldownPeriod string `json:"cooldownPeriod,omitempty"`
}

// JAXStatsAlertStatus defines the observed state of JAXStatsAlert.
type JAXStatsAlertStatus struct {
	// LastEvaluationTime is the timestamp of the last alert evaluation
	LastEvaluationTime *metav1.Time `json:"lastEvaluationTime,omitempty"`

	// AlertStatus indicates the current status of the alert
	AlertStatus string `json:"alertStatus,omitempty"`

	// ActiveAlerts tracks currently active/firing alert rule names
	ActiveAlerts []string `json:"activeAlerts,omitempty"`

	// FiringAlertsCount tracks the number of currently firing alerts
	FiringAlertsCount int32 `json:"firingAlertsCount,omitempty"`

	// LastFiredTime is the timestamp of the last time an alert fired
	LastFiredTime *metav1.Time `json:"lastFiredTime,omitempty"`

	// ErrorMessage contains any error message if alerting failed
	ErrorMessage string `json:"errorMessage,omitempty"`

	// Conditions represent the latest available observations of the alert's current state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=jsa
// +kubebuilder:printcolumn:name="Enabled",type=boolean,JSONPath=`.spec.enabled`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.alertStatus`
// +kubebuilder:printcolumn:name="Firing",type=integer,JSONPath=`.status.firingAlertsCount`
// +kubebuilder:printcolumn:name="Collector",type=string,JSONPath=`.spec.collectorRef`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// JAXStatsAlert is the Schema for the jaxstatsalerts API
type JAXStatsAlert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JAXStatsAlertSpec   `json:"spec,omitempty"`
	Status JAXStatsAlertStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JAXStatsAlertList contains a list of JAXStatsAlert
type JAXStatsAlertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JAXStatsAlert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JAXStatsAlert{}, &JAXStatsAlertList{})
}
