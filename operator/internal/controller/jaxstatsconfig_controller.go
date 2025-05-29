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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	statsv1alpha1 "github.com/raph/corium/operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JAXStatsConfigReconciler reconciles a JAXStatsConfig object
type JAXStatsConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JAXStatsConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *JAXStatsConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the JAXStatsConfig instance
	jaxStatsConfig := &statsv1alpha1.JAXStatsConfig{}
	if err := r.Get(ctx, req.NamespacedName, jaxStatsConfig); err != nil {
		log.Error(err, "unable to fetch JAXStatsConfig")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Update status with current time
	now := metav1.Now()
	jaxStatsConfig.Status.LastCollectionTime = &now

	// Set collection status based on enabled flag
	if jaxStatsConfig.Spec.Enabled {
		jaxStatsConfig.Status.CollectionStatus = "Active"
	} else {
		jaxStatsConfig.Status.CollectionStatus = "Disabled"
	}

	// Add a condition to track the configuration state
	condition := metav1.Condition{
		Type:               "ConfigurationValid",
		Status:             metav1.ConditionTrue,
		Reason:             "ConfigurationValid",
		Message:            "Configuration is valid and being applied",
		LastTransitionTime: now,
	}

	// Validate storage configuration
	if jaxStatsConfig.Spec.StorageConfig.Type == "" {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "InvalidStorageConfig"
		condition.Message = "Storage type is required"
		jaxStatsConfig.Status.ErrorMessage = "Storage type is required"
	}

	// Update the status
	if err := r.Status().Update(ctx, jaxStatsConfig); err != nil {
		log.Error(err, "unable to update JAXStatsConfig status")
		return ctrl.Result{}, err
	}

	// Set up requeue interval based on collection interval
	requeueAfter := time.Duration(jaxStatsConfig.Spec.CollectionInterval) * time.Second
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JAXStatsConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&statsv1alpha1.JAXStatsConfig{}).
		Named("jaxstatsconfig").
		Complete(r)
}
