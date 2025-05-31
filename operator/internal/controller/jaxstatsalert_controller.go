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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	statsv1alpha1 "github.com/raph/corium/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JAXStatsAlertReconciler reconciles a JAXStatsAlert object
type JAXStatsAlertReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsalerts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsalerts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsalerts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JAXStatsAlert object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *JAXStatsAlertReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the JAXStatsAlert instance
	alert := &statsv1alpha1.JAXStatsAlert{}
	if err := r.Get(ctx, req.NamespacedName, alert); err != nil {
		logger.Error(err, "unable to fetch JAXStatsAlert")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the referenced JAXStatsCollector
	collector := &statsv1alpha1.JAXStatsCollector{}
	if err := r.Get(ctx, client.ObjectKey{
		Namespace: req.Namespace,
		Name:      alert.Spec.CollectorRef,
	}, collector); err != nil {
		logger.Error(err, "unable to fetch referenced JAXStatsCollector")
		alert.Status.ErrorMessage = fmt.Sprintf("Failed to fetch collector: %v", err)
		alert.Status.AlertStatus = "Error"
		if err := r.Status().Update(ctx, alert); err != nil {
			logger.Error(err, "unable to update alert status")
		}
		return ctrl.Result{}, err
	}

	// Update status
	now := metav1.Now()
	alert.Status.LastAlertTime = &now

	// Validate alert configuration
	condition := metav1.Condition{
		Type:               "AlertValid",
		Status:             metav1.ConditionTrue,
		Reason:             "AlertValid",
		Message:            "Alert configuration is valid",
		LastTransitionTime: now,
	}

	// Check if alert is enabled
	if !alert.Spec.Enabled {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "AlertDisabled"
		condition.Message = "Alert is disabled"
		alert.Status.AlertStatus = "Disabled"
	} else {
		alert.Status.AlertStatus = "Active"
	}

	// Validate rules
	for _, rule := range alert.Spec.Rules {
		if rule.Metric == "" || rule.Operator == "" {
			condition.Status = metav1.ConditionFalse
			condition.Reason = "InvalidRule"
			condition.Message = fmt.Sprintf("Invalid rule configuration: %s", rule.Name)
			alert.Status.ErrorMessage = condition.Message
			break
		}
	}

	// Update the status
	if err := r.Status().Update(ctx, alert); err != nil {
		logger.Error(err, "unable to update JAXStatsAlert status")
		return ctrl.Result{}, err
	}

	// Calculate next reconciliation time based on cooldown period
	// For now, we'll use a simple 1-minute interval
	requeueAfter := time.Minute
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JAXStatsAlertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&statsv1alpha1.JAXStatsAlert{}).
		Named("jaxstatsalert").
		Complete(r)
}
