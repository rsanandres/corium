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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	statsv1alpha1 "github.com/raph/corium/operator/api/v1alpha1"
)

var _ = Describe("JAXStatsAlert Controller", func() {
	const namespace = "default"

	createMetricsConfigMap := func(ctx context.Context, name, ns string, metrics CollectedMetrics) {
		data, err := json.MarshalIndent(metrics, "", "  ")
		Expect(err).NotTo(HaveOccurred())
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
			Data:       map[string]string{"metrics.json": string(data)},
		}
		existing := &corev1.ConfigMap{}
		if err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, existing); err == nil {
			existing.Data = cm.Data
			Expect(k8sClient.Update(ctx, existing)).To(Succeed())
		} else {
			Expect(k8sClient.Create(ctx, cm)).To(Succeed())
		}
	}

	Context("When an alert fires due to high restart count", func() {
		const (
			collectorName = "alert-test-collector"
			alertName     = "test-alert-firing"
			cmName        = "alert-test-collector-metrics"
		)
		ctx := context.Background()
		alertNN := types.NamespacedName{Name: alertName, Namespace: namespace}

		BeforeEach(func() {
			// Create collector
			collector := &statsv1alpha1.JAXStatsCollector{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: collectorName, Namespace: namespace}, collector)
			if err != nil && errors.IsNotFound(err) {
				collector = &statsv1alpha1.JAXStatsCollector{
					ObjectMeta: metav1.ObjectMeta{Name: collectorName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsCollectorSpec{
						TargetNamespace: namespace,
						ConfigRef:       "some-config",
						Metrics:         []string{"restart_count"},
						Selector: metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "test"},
						},
					},
				}
				Expect(k8sClient.Create(ctx, collector)).To(Succeed())
			}
			// Set MetricsConfigMap in status
			collector = &statsv1alpha1.JAXStatsCollector{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: collectorName, Namespace: namespace}, collector)).To(Succeed())
			collector.Status.MetricsConfigMap = cmName
			k8sClient.Status().Update(ctx, collector)

			// Create metrics ConfigMap with high restart count
			createMetricsConfigMap(ctx, cmName, namespace, CollectedMetrics{
				CollectorName: collectorName,
				PodCount:      2,
				Pods: []PodMetrics{
					{Name: "pod-1", RestartCount: 10, Ready: true, ContainerCount: 1},
					{Name: "pod-2", RestartCount: 5, Ready: false, ContainerCount: 1},
				},
				Summary: MetricsSummary{
					TotalPods:     2,
					ReadyPods:     1,
					NotReadyPods:  1,
					TotalRestarts: 15,
				},
			})

			// Create alert
			alert := &statsv1alpha1.JAXStatsAlert{}
			err = k8sClient.Get(ctx, alertNN, alert)
			if err != nil && errors.IsNotFound(err) {
				alert = &statsv1alpha1.JAXStatsAlert{
					ObjectMeta: metav1.ObjectMeta{Name: alertName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsAlertSpec{
						Enabled:        true,
						CollectorRef:   collectorName,
						CooldownPeriod: "1m",
						Rules: []statsv1alpha1.AlertRule{
							{
								Name:      "high-restarts",
								Metric:    "restart_count",
								Operator:  ">",
								Threshold: "5",
								Severity:  "warning",
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, alert)).To(Succeed())
			}
		})

		AfterEach(func() {
			alert := &statsv1alpha1.JAXStatsAlert{}
			if err := k8sClient.Get(ctx, alertNN, alert); err == nil {
				k8sClient.Delete(ctx, alert)
			}
			cm := &corev1.ConfigMap{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: cmName, Namespace: namespace}, cm); err == nil {
				k8sClient.Delete(ctx, cm)
			}
			collector := &statsv1alpha1.JAXStatsCollector{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: collectorName, Namespace: namespace}, collector); err == nil {
				k8sClient.Delete(ctx, collector)
			}
		})

		It("should set alert status to Firing with active alerts", func() {
			fakeRecorder := record.NewFakeRecorder(10)
			reconciler := &JAXStatsAlertReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: fakeRecorder,
			}

			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: alertNN})
			Expect(err).NotTo(HaveOccurred())

			updated := &statsv1alpha1.JAXStatsAlert{}
			Expect(k8sClient.Get(ctx, alertNN, updated)).To(Succeed())
			Expect(updated.Status.AlertStatus).To(Equal("Firing"))
			Expect(updated.Status.FiringAlertsCount).To(Equal(int32(1)))
			Expect(updated.Status.ActiveAlerts).To(ContainElement("high-restarts"))

			// Verify an event was emitted
			Expect(fakeRecorder.Events).To(HaveLen(1))
		})
	})

	Context("When alert resolves (below threshold)", func() {
		const (
			collectorName = "alert-resolve-collector"
			alertName     = "test-alert-resolving"
			cmName        = "alert-resolve-collector-metrics"
		)
		ctx := context.Background()
		alertNN := types.NamespacedName{Name: alertName, Namespace: namespace}

		BeforeEach(func() {
			collector := &statsv1alpha1.JAXStatsCollector{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: collectorName, Namespace: namespace}, collector)
			if err != nil && errors.IsNotFound(err) {
				collector = &statsv1alpha1.JAXStatsCollector{
					ObjectMeta: metav1.ObjectMeta{Name: collectorName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsCollectorSpec{
						TargetNamespace: namespace,
						ConfigRef:       "some-config",
						Metrics:         []string{"restart_count"},
						Selector: metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "test"},
						},
					},
				}
				Expect(k8sClient.Create(ctx, collector)).To(Succeed())
			}
			collector = &statsv1alpha1.JAXStatsCollector{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: collectorName, Namespace: namespace}, collector)).To(Succeed())
			collector.Status.MetricsConfigMap = cmName
			k8sClient.Status().Update(ctx, collector)

			// Low restart count — should NOT fire
			createMetricsConfigMap(ctx, cmName, namespace, CollectedMetrics{
				CollectorName: collectorName,
				PodCount:      1,
				Pods: []PodMetrics{
					{Name: "pod-1", RestartCount: 2, Ready: true, ContainerCount: 1},
				},
				Summary: MetricsSummary{
					TotalPods:     1,
					ReadyPods:     1,
					TotalRestarts: 2,
				},
			})

			alert := &statsv1alpha1.JAXStatsAlert{}
			err = k8sClient.Get(ctx, alertNN, alert)
			if err != nil && errors.IsNotFound(err) {
				alert = &statsv1alpha1.JAXStatsAlert{
					ObjectMeta: metav1.ObjectMeta{Name: alertName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsAlertSpec{
						Enabled:        true,
						CollectorRef:   collectorName,
						CooldownPeriod: "1m",
						Rules: []statsv1alpha1.AlertRule{
							{
								Name:      "high-restarts",
								Metric:    "restart_count",
								Operator:  ">",
								Threshold: "5",
								Severity:  "warning",
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, alert)).To(Succeed())
			}
		})

		AfterEach(func() {
			alert := &statsv1alpha1.JAXStatsAlert{}
			if err := k8sClient.Get(ctx, alertNN, alert); err == nil {
				k8sClient.Delete(ctx, alert)
			}
			cm := &corev1.ConfigMap{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: cmName, Namespace: namespace}, cm); err == nil {
				k8sClient.Delete(ctx, cm)
			}
			collector := &statsv1alpha1.JAXStatsCollector{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: collectorName, Namespace: namespace}, collector); err == nil {
				k8sClient.Delete(ctx, collector)
			}
		})

		It("should set alert status to OK with no active alerts", func() {
			reconciler := &JAXStatsAlertReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(10),
			}

			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: alertNN})
			Expect(err).NotTo(HaveOccurred())

			updated := &statsv1alpha1.JAXStatsAlert{}
			Expect(k8sClient.Get(ctx, alertNN, updated)).To(Succeed())
			Expect(updated.Status.AlertStatus).To(Equal("OK"))
			Expect(updated.Status.FiringAlertsCount).To(Equal(int32(0)))
			Expect(updated.Status.ActiveAlerts).To(BeEmpty())
		})
	})

	Context("When alert is disabled", func() {
		const alertName = "test-alert-disabled"
		ctx := context.Background()
		alertNN := types.NamespacedName{Name: alertName, Namespace: namespace}

		BeforeEach(func() {
			alert := &statsv1alpha1.JAXStatsAlert{}
			err := k8sClient.Get(ctx, alertNN, alert)
			if err != nil && errors.IsNotFound(err) {
				alert = &statsv1alpha1.JAXStatsAlert{
					ObjectMeta: metav1.ObjectMeta{Name: alertName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsAlertSpec{
						Enabled:      false,
						CollectorRef: "some-collector",
						Rules: []statsv1alpha1.AlertRule{
							{
								Name:      "test-rule",
								Metric:    "restart_count",
								Operator:  ">",
								Threshold: "5",
								Severity:  "warning",
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, alert)).To(Succeed())
			}
		})

		AfterEach(func() {
			alert := &statsv1alpha1.JAXStatsAlert{}
			if err := k8sClient.Get(ctx, alertNN, alert); err == nil {
				k8sClient.Delete(ctx, alert)
			}
		})

		It("should set status to Disabled", func() {
			reconciler := &JAXStatsAlertReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(10),
			}

			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: alertNN})
			Expect(err).NotTo(HaveOccurred())

			updated := &statsv1alpha1.JAXStatsAlert{}
			Expect(k8sClient.Get(ctx, alertNN, updated)).To(Succeed())
			Expect(updated.Status.AlertStatus).To(Equal("Disabled"))
			Expect(updated.Status.FiringAlertsCount).To(Equal(int32(0)))
		})
	})
})
