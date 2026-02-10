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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	statsv1alpha1 "github.com/raph/corium/operator/api/v1alpha1"
)

// JAXStatsAlertReconciler reconciles a JAXStatsAlert object
type JAXStatsAlertReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsalerts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsalerts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsalerts/finalizers,verbs=update
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatscollectors,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *JAXStatsAlertReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	alert := &statsv1alpha1.JAXStatsAlert{}
	if err := r.Get(ctx, req.NamespacedName, alert); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	now := metav1.Now()
	alert.Status.LastEvaluationTime = &now

	// Check if alert is enabled
	if !alert.Spec.Enabled {
		alert.Status.AlertStatus = "Disabled"
		alert.Status.ActiveAlerts = nil
		alert.Status.FiringAlertsCount = 0
		alert.Status.ErrorMessage = ""
		meta.SetStatusCondition(&alert.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			Reason:             "AlertDisabled",
			Message:            "Alert is disabled",
			LastTransitionTime: now,
			ObservedGeneration: alert.Generation,
		})
		if err := r.Status().Update(ctx, alert); err != nil {
			return ctrl.Result{}, err
		}
		activeAlertsGauge.WithLabelValues(alert.Name).Set(0)
		return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
	}

	// Fetch the referenced collector
	collector := &statsv1alpha1.JAXStatsCollector{}
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      alert.Spec.CollectorRef,
	}, collector); err != nil {
		reconcileErrorsTotal.WithLabelValues("alert").Inc()
		alert.Status.ErrorMessage = fmt.Sprintf("collector %q not found", alert.Spec.CollectorRef)
		alert.Status.AlertStatus = "Error"
		meta.SetStatusCondition(&alert.Status.Conditions, metav1.Condition{
			Type:               "CollectorAvailable",
			Status:             metav1.ConditionFalse,
			Reason:             "CollectorNotFound",
			Message:            alert.Status.ErrorMessage,
			LastTransitionTime: now,
			ObservedGeneration: alert.Generation,
		})
		r.Status().Update(ctx, alert)
		r.Recorder.Event(alert, "Warning", "CollectorNotFound", alert.Status.ErrorMessage)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	meta.SetStatusCondition(&alert.Status.Conditions, metav1.Condition{
		Type:               "CollectorAvailable",
		Status:             metav1.ConditionTrue,
		Reason:             "CollectorReady",
		Message:            "Referenced collector is available",
		LastTransitionTime: now,
		ObservedGeneration: alert.Generation,
	})

	// Read metrics from the collector's ConfigMap
	cmName := collector.Status.MetricsConfigMap
	if cmName == "" {
		cmName = collector.Name + "-metrics"
	}

	cm := &corev1.ConfigMap{}
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      cmName,
	}, cm); err != nil {
		reconcileErrorsTotal.WithLabelValues("alert").Inc()
		alert.Status.ErrorMessage = fmt.Sprintf("metrics ConfigMap %q not found", cmName)
		alert.Status.AlertStatus = "Pending"
		r.Status().Update(ctx, alert)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Parse collected metrics
	metricsJSON, ok := cm.Data["metrics.json"]
	if !ok {
		alert.Status.ErrorMessage = "metrics.json key not found in ConfigMap"
		alert.Status.AlertStatus = "Error"
		r.Status().Update(ctx, alert)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	var collectedMetrics CollectedMetrics
	if err := json.Unmarshal([]byte(metricsJSON), &collectedMetrics); err != nil {
		reconcileErrorsTotal.WithLabelValues("alert").Inc()
		alert.Status.ErrorMessage = fmt.Sprintf("failed to parse metrics: %v", err)
		alert.Status.AlertStatus = "Error"
		r.Status().Update(ctx, alert)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Evaluate alert rules
	evalResults := EvaluateAlertRules(alert.Spec.Rules, &collectedMetrics)

	// Process results: determine firing/resolved alerts
	var firingAlerts []string
	for _, result := range evalResults {
		if result.Firing {
			firingAlerts = append(firingAlerts, result.RuleName)
		}
	}

	// Check cooldown before emitting events
	cooldown := parseDuration(alert.Spec.CooldownPeriod, 5*time.Minute)
	canEmit := alert.Status.LastFiredTime == nil ||
		time.Since(alert.Status.LastFiredTime.Time) > cooldown

	// Emit events for firing alerts
	previousFiring := make(map[string]bool)
	for _, a := range alert.Status.ActiveAlerts {
		previousFiring[a] = true
	}

	for _, result := range evalResults {
		if result.Firing && canEmit {
			eventType := "Warning"
			if result.Severity == "critical" {
				eventType = "Warning"
			}
			r.Recorder.Event(alert, eventType, "AlertFiring",
				fmt.Sprintf("[%s] %s: %s=%.0f exceeds threshold %.0f",
					result.Severity, result.RuleName, result.Metric, result.ActualValue, result.Threshold))
			logger.Info("alert firing", "rule", result.RuleName, "metric", result.Metric,
				"actual", result.ActualValue, "threshold", result.Threshold)
		}
		if !result.Firing && previousFiring[result.RuleName] {
			r.Recorder.Event(alert, "Normal", "AlertResolved",
				fmt.Sprintf("Resolved: %s (%s=%.0f, threshold=%.0f)",
					result.RuleName, result.Metric, result.ActualValue, result.Threshold))
			logger.Info("alert resolved", "rule", result.RuleName)
		}
	}

	// Update status
	alert.Status.ActiveAlerts = firingAlerts
	alert.Status.FiringAlertsCount = int32(len(firingAlerts))
	alert.Status.ErrorMessage = ""

	if len(firingAlerts) > 0 {
		alert.Status.AlertStatus = "Firing"
		if canEmit {
			alert.Status.LastFiredTime = &now
		}
		meta.SetStatusCondition(&alert.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			Reason:             "AlertsFiring",
			Message:            fmt.Sprintf("%d alert(s) firing", len(firingAlerts)),
			LastTransitionTime: now,
			ObservedGeneration: alert.Generation,
		})
	} else {
		alert.Status.AlertStatus = "OK"
		meta.SetStatusCondition(&alert.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			Reason:             "AllClear",
			Message:            "No alerts firing",
			LastTransitionTime: now,
			ObservedGeneration: alert.Generation,
		})
	}

	if err := r.Status().Update(ctx, alert); err != nil {
		logger.Error(err, "unable to update JAXStatsAlert status")
		return ctrl.Result{}, err
	}

	// Update Prometheus metric
	activeAlertsGauge.WithLabelValues(alert.Name).Set(float64(len(firingAlerts)))

	return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
}

func parseDuration(s string, fallback time.Duration) time.Duration {
	if s == "" {
		return fallback
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fallback
	}
	return d
}

// SetupWithManager sets up the controller with the Manager.
func (r *JAXStatsAlertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&statsv1alpha1.JAXStatsAlert{}).
		Named("jaxstatsalert").
		Complete(r)
}
