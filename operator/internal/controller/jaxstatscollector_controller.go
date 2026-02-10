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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	statsv1alpha1 "github.com/raph/corium/operator/api/v1alpha1"
)

// JAXStatsCollectorReconciler reconciles a JAXStatsCollector object
type JAXStatsCollectorReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatscollectors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatscollectors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatscollectors/finalizers,verbs=update
// +kubebuilder:rbac:groups=stats.corium.io,resources=jaxstatsconfigs,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *JAXStatsCollectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	collector := &statsv1alpha1.JAXStatsCollector{}
	if err := r.Get(ctx, req.NamespacedName, collector); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch the referenced JAXStatsConfig
	config := &statsv1alpha1.JAXStatsConfig{}
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      collector.Spec.ConfigRef,
	}, config); err != nil {
		reconcileErrorsTotal.WithLabelValues("collector").Inc()
		collector.Status.ErrorMessage = fmt.Sprintf("config %q not found: %v", collector.Spec.ConfigRef, err)
		collector.Status.CollectionStatus = "Error"
		meta.SetStatusCondition(&collector.Status.Conditions, metav1.Condition{
			Type:               "ConfigAvailable",
			Status:             metav1.ConditionFalse,
			Reason:             "ConfigNotFound",
			Message:            collector.Status.ErrorMessage,
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: collector.Generation,
		})
		if updateErr := r.Status().Update(ctx, collector); updateErr != nil {
			logger.Error(updateErr, "unable to update collector status")
		}
		r.Recorder.Event(collector, "Warning", "ConfigNotFound", collector.Status.ErrorMessage)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Check if config is enabled
	if !config.Spec.Enabled {
		collector.Status.CollectionStatus = "Disabled"
		collector.Status.ErrorMessage = ""
		meta.SetStatusCondition(&collector.Status.Conditions, metav1.Condition{
			Type:               "ConfigAvailable",
			Status:             metav1.ConditionFalse,
			Reason:             "ConfigDisabled",
			Message:            "Referenced JAXStatsConfig is disabled",
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: collector.Generation,
		})
		if err := r.Status().Update(ctx, collector); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
	}

	meta.SetStatusCondition(&collector.Status.Conditions, metav1.Condition{
		Type:               "ConfigAvailable",
		Status:             metav1.ConditionTrue,
		Reason:             "ConfigReady",
		Message:            "Referenced JAXStatsConfig is available and enabled",
		LastTransitionTime: metav1.Now(),
		ObservedGeneration: collector.Generation,
	})

	// Discover pods by label selector
	selector, err := metav1.LabelSelectorAsSelector(&collector.Spec.Selector)
	if err != nil {
		reconcileErrorsTotal.WithLabelValues("collector").Inc()
		collector.Status.ErrorMessage = fmt.Sprintf("invalid label selector: %v", err)
		collector.Status.CollectionStatus = "Error"
		r.Status().Update(ctx, collector)
		return ctrl.Result{}, err
	}

	podList := &corev1.PodList{}
	listOpts := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     collector.Spec.TargetNamespace,
	}
	if err := r.List(ctx, podList, listOpts); err != nil {
		reconcileErrorsTotal.WithLabelValues("collector").Inc()
		logger.Error(err, "unable to list pods", "namespace", collector.Spec.TargetNamespace, "selector", selector.String())
		collector.Status.ErrorMessage = fmt.Sprintf("failed to list pods: %v", err)
		collector.Status.CollectionStatus = "Error"
		r.Status().Update(ctx, collector)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Collect metrics from discovered pods
	podMetrics := CollectPodMetrics(podList.Items)
	collectedMetrics := BuildCollectedMetrics(collector.Name, collector.Spec.TargetNamespace, podMetrics)

	metricsJSON, err := MarshalMetrics(collectedMetrics)
	if err != nil {
		reconcileErrorsTotal.WithLabelValues("collector").Inc()
		return ctrl.Result{}, fmt.Errorf("failed to marshal metrics: %w", err)
	}

	// Create or update the metrics ConfigMap
	cmName := collector.Name + "-metrics"
	if err := r.ensureMetricsConfigMap(ctx, collector, cmName, metricsJSON); err != nil {
		reconcileErrorsTotal.WithLabelValues("collector").Inc()
		logger.Error(err, "unable to ensure metrics ConfigMap")
		collector.Status.ErrorMessage = fmt.Sprintf("failed to manage ConfigMap: %v", err)
		collector.Status.CollectionStatus = "Error"
		r.Status().Update(ctx, collector)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Update status
	now := metav1.Now()
	podCount := int32(len(podList.Items))
	collector.Status.LastCollectionTime = &now
	collector.Status.CollectionStatus = "Active"
	collector.Status.DiscoveredPods = podCount
	collector.Status.CollectedResources = podCount
	collector.Status.MetricsConfigMap = cmName
	collector.Status.ErrorMessage = ""

	meta.SetStatusCondition(&collector.Status.Conditions, metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             "CollectionSuccessful",
		Message:            fmt.Sprintf("Discovered %d pods, metrics stored in ConfigMap %s", podCount, cmName),
		LastTransitionTime: now,
		ObservedGeneration: collector.Generation,
	})

	if err := r.Status().Update(ctx, collector); err != nil {
		logger.Error(err, "unable to update JAXStatsCollector status")
		return ctrl.Result{}, err
	}

	// Update Prometheus metric
	discoveredPodsGauge.WithLabelValues(collector.Name).Set(float64(podCount))

	logger.Info("collection complete", "pods", podCount, "configmap", cmName)

	requeueAfter := time.Duration(config.Spec.CollectionInterval) * time.Second
	if requeueAfter == 0 {
		requeueAfter = 60 * time.Second
	}
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

func (r *JAXStatsCollectorReconciler) ensureMetricsConfigMap(
	ctx context.Context,
	collector *statsv1alpha1.JAXStatsCollector,
	cmName string,
	metricsJSON string,
) error {
	cm := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: collector.Namespace,
		Name:      cmName,
	}, cm)

	if errors.IsNotFound(err) {
		cm = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: collector.Namespace,
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "jaxstats-operator",
					"stats.corium.io/collector":    collector.Name,
				},
			},
			Data: map[string]string{
				"metrics.json": metricsJSON,
			},
		}
		if err := controllerutil.SetControllerReference(collector, cm, r.Scheme); err != nil {
			return fmt.Errorf("failed to set owner reference: %w", err)
		}
		if err := r.Create(ctx, cm); err != nil {
			return fmt.Errorf("failed to create ConfigMap: %w", err)
		}
		r.Recorder.Event(collector, "Normal", "ConfigMapCreated", fmt.Sprintf("Created metrics ConfigMap %s", cmName))
		return nil
	}
	if err != nil {
		return err
	}

	// Update existing ConfigMap
	cm.Data["metrics.json"] = metricsJSON
	if err := r.Update(ctx, cm); err != nil {
		return fmt.Errorf("failed to update ConfigMap: %w", err)
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JAXStatsCollectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&statsv1alpha1.JAXStatsCollector{}).
		Owns(&corev1.ConfigMap{}).
		Named("jaxstatscollector").
		Complete(r)
}
