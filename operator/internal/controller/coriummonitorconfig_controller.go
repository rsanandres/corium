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
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	monitorv1alpha1 "github.com/raph/corium/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	configFinalizerName = "monitor.corium.io/config-cleanup"
)

// CoriumMonitorConfigReconciler reconciles a CoriumMonitorConfig object
type CoriumMonitorConfigReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=monitor.corium.io,resources=coriummonitorconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitor.corium.io,resources=coriummonitorconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitor.corium.io,resources=coriummonitorconfigs/finalizers,verbs=update
// +kubebuilder:rbac:groups=monitor.corium.io,resources=coriummonitorcollectors,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *CoriumMonitorConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	config := &monitorv1alpha1.CoriumMonitorConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle deletion with finalizer
	if !config.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(config, configFinalizerName) {
			logger.Info("running finalizer for CoriumMonitorConfig")
			if err := r.cleanupDependents(ctx, config); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(config, configFinalizerName)
			if err := r.Update(ctx, config); err != nil {
				return ctrl.Result{}, err
			}
			r.Recorder.Event(config, "Normal", "Deleted", "Config and dependent resources cleaned up")
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(config, configFinalizerName) {
		controllerutil.AddFinalizer(config, configFinalizerName)
		if err := r.Update(ctx, config); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Validate spec
	if err := r.validateConfig(ctx, config); err != nil {
		reconcileErrorsTotal.WithLabelValues("config").Inc()
		return ctrl.Result{}, err
	}

	// Count dependent collectors
	collectors := &monitorv1alpha1.CoriumMonitorCollectorList{}
	if err := r.List(ctx, collectors, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "unable to list collectors")
		return ctrl.Result{}, err
	}
	var depCount int32
	for _, col := range collectors.Items {
		if col.Spec.ConfigRef == config.Name {
			depCount++
		}
	}
	config.Status.DependentCollectors = depCount

	// Set collection status
	now := metav1.Now()
	config.Status.LastCollectionTime = &now
	if config.Spec.Enabled {
		config.Status.CollectionStatus = "Active"
	} else {
		config.Status.CollectionStatus = "Disabled"
		r.Recorder.Event(config, "Normal", "Disabled", "Config collection is disabled")
	}

	if err := r.Status().Update(ctx, config); err != nil {
		logger.Error(err, "unable to update CoriumMonitorConfig status")
		return ctrl.Result{}, err
	}

	requeueAfter := time.Duration(config.Spec.CollectionInterval) * time.Second
	if requeueAfter == 0 {
		requeueAfter = 60 * time.Second
	}
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

func (r *CoriumMonitorConfigReconciler) validateConfig(ctx context.Context, config *monitorv1alpha1.CoriumMonitorConfig) error {
	now := metav1.Now()

	// Validate storage configuration
	validStorage := config.Spec.StorageConfig.Type != ""
	condition := metav1.Condition{
		Type:               "ConfigurationValid",
		LastTransitionTime: now,
		ObservedGeneration: config.Generation,
	}

	if !validStorage {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "InvalidStorageConfig"
		condition.Message = "Storage type is required"
		config.Status.ErrorMessage = "Storage type is required"
		r.Recorder.Event(config, "Warning", "ValidationFailed", "Storage type is required")
	} else {
		condition.Status = metav1.ConditionTrue
		condition.Reason = "ConfigurationValid"
		condition.Message = "Configuration is valid"
		config.Status.ErrorMessage = ""
	}
	meta.SetStatusCondition(&config.Status.Conditions, condition)

	return r.Status().Update(ctx, config)
}

func (r *CoriumMonitorConfigReconciler) cleanupDependents(ctx context.Context, config *monitorv1alpha1.CoriumMonitorConfig) error {
	logger := log.FromContext(ctx)

	collectors := &monitorv1alpha1.CoriumMonitorCollectorList{}
	if err := r.List(ctx, collectors, client.InNamespace(config.Namespace)); err != nil {
		return err
	}

	for i := range collectors.Items {
		col := &collectors.Items[i]
		if col.Spec.ConfigRef == config.Name {
			logger.Info("marking dependent collector as orphaned", "collector", col.Name)
			meta.SetStatusCondition(&col.Status.Conditions, metav1.Condition{
				Type:               "ConfigAvailable",
				Status:             metav1.ConditionFalse,
				Reason:             "ConfigDeleted",
				Message:            "Referenced CoriumMonitorConfig has been deleted",
				LastTransitionTime: metav1.Now(),
			})
			col.Status.CollectionStatus = alertStatusError
			col.Status.ErrorMessage = "Referenced config was deleted"
			if err := r.Status().Update(ctx, col); err != nil {
				logger.Error(err, "unable to update orphaned collector status", "collector", col.Name)
			}
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CoriumMonitorConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitorv1alpha1.CoriumMonitorConfig{}).
		Named("coriummonitorconfig").
		Complete(r)
}
