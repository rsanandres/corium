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

// JAXStatsCollectorReconciler reconciles a JAXStatsCollector object
type JAXStatsCollectorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatscollectors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatscollectors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatscollectors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JAXStatsCollector object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *JAXStatsCollectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the JAXStatsCollector instance
	collector := &statsv1alpha1.JAXStatsCollector{}
	if err := r.Get(ctx, req.NamespacedName, collector); err != nil {
		log.Error(err, "unable to fetch JAXStatsCollector")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the referenced JAXStatsConfig
	config := &statsv1alpha1.JAXStatsConfig{}
	if err := r.Get(ctx, client.ObjectKey{
		Namespace: req.Namespace,
		Name:      collector.Spec.ConfigRef,
	}, config); err != nil {
		log.Error(err, "unable to fetch referenced JAXStatsConfig")
		collector.Status.ErrorMessage = fmt.Sprintf("Failed to fetch config: %v", err)
		collector.Status.CollectionStatus = "Error"
		if err := r.Status().Update(ctx, collector); err != nil {
			log.Error(err, "unable to update collector status")
		}
		return ctrl.Result{}, err
	}

	// Update status
	now := metav1.Now()
	collector.Status.LastCollectionTime = &now

	// Validate collector configuration
	condition := metav1.Condition{
		Type:               "CollectorValid",
		Status:             metav1.ConditionTrue,
		Reason:             "CollectorValid",
		Message:            "Collector configuration is valid",
		LastTransitionTime: now,
	}

	// Check if config is enabled
	if !config.Spec.Enabled {
		condition.Status = metav1.ConditionFalse
		condition.Reason = "ConfigDisabled"
		condition.Message = "Referenced JAXStatsConfig is disabled"
		collector.Status.CollectionStatus = "Disabled"
	} else {
		collector.Status.CollectionStatus = "Active"
	}

	// Update the status
	if err := r.Status().Update(ctx, collector); err != nil {
		log.Error(err, "unable to update JAXStatsCollector status")
		return ctrl.Result{}, err
	}

	// Calculate next reconciliation time based on collection schedule
	// For now, we'll use a simple 5-minute interval
	requeueAfter := 5 * time.Minute
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JAXStatsCollectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&statsv1alpha1.JAXStatsCollector{}).
		Named("jaxstatscollector").
		Complete(r)
}
