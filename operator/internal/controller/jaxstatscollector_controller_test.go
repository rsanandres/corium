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

var _ = Describe("JAXStatsCollector Controller", func() {
	const namespace = "default"

	Context("When reconciling with a valid config and matching pods", func() {
		const (
			configName    = "collector-test-config"
			collectorName = "test-collector-pods"
			podName       = "test-pod-1"
		)
		ctx := context.Background()
		collectorNN := types.NamespacedName{Name: collectorName, Namespace: namespace}

		BeforeEach(func() {
			// Create config
			config := &statsv1alpha1.JAXStatsConfig{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: configName, Namespace: namespace}, config)
			if err != nil && errors.IsNotFound(err) {
				config = &statsv1alpha1.JAXStatsConfig{
					ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsConfigSpec{
						Enabled:            true,
						CollectionInterval: 30,
						Metrics:            []string{"memory_usage"},
						StorageConfig:      statsv1alpha1.StorageConfig{Type: "configmap"},
					},
				}
				Expect(k8sClient.Create(ctx, config)).To(Succeed())
			}

			// Create a pod with matching labels
			pod := &corev1.Pod{}
			err = k8sClient.Get(ctx, types.NamespacedName{Name: podName, Namespace: namespace}, pod)
			if err != nil && errors.IsNotFound(err) {
				pod = &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      podName,
						Namespace: namespace,
						Labels:    map[string]string{"app": "jax-test"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{Name: "main", Image: "nginx:latest"},
						},
					},
				}
				Expect(k8sClient.Create(ctx, pod)).To(Succeed())
			}

			// Create collector
			collector := &statsv1alpha1.JAXStatsCollector{}
			err = k8sClient.Get(ctx, collectorNN, collector)
			if err != nil && errors.IsNotFound(err) {
				collector = &statsv1alpha1.JAXStatsCollector{
					ObjectMeta: metav1.ObjectMeta{Name: collectorName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsCollectorSpec{
						TargetNamespace: namespace,
						ConfigRef:       configName,
						Metrics:         []string{"memory_usage"},
						Selector: metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "jax-test"},
						},
					},
				}
				Expect(k8sClient.Create(ctx, collector)).To(Succeed())
			}
		})

		AfterEach(func() {
			// Cleanup ConfigMap
			cm := &corev1.ConfigMap{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: collectorName + "-metrics", Namespace: namespace}, cm); err == nil {
				k8sClient.Delete(ctx, cm)
			}
			// Cleanup collector
			collector := &statsv1alpha1.JAXStatsCollector{}
			if err := k8sClient.Get(ctx, collectorNN, collector); err == nil {
				k8sClient.Delete(ctx, collector)
			}
			// Cleanup pod
			pod := &corev1.Pod{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: podName, Namespace: namespace}, pod); err == nil {
				k8sClient.Delete(ctx, pod)
			}
			// Cleanup config
			config := &statsv1alpha1.JAXStatsConfig{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: configName, Namespace: namespace}, config); err == nil {
				config.Finalizers = nil
				k8sClient.Update(ctx, config)
				k8sClient.Delete(ctx, config)
			}
		})

		It("should discover pods and create a metrics ConfigMap", func() {
			reconciler := &JAXStatsCollectorReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(10),
			}

			result, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: collectorNN})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeNumerically(">", 0))

			// Verify collector status was updated
			updated := &statsv1alpha1.JAXStatsCollector{}
			Expect(k8sClient.Get(ctx, collectorNN, updated)).To(Succeed())
			Expect(updated.Status.CollectionStatus).To(Equal("Active"))
			Expect(updated.Status.DiscoveredPods).To(BeNumerically(">=", 1))
			Expect(updated.Status.MetricsConfigMap).To(Equal(collectorName + "-metrics"))

			// Verify ConfigMap was created
			cm := &corev1.ConfigMap{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      collectorName + "-metrics",
				Namespace: namespace,
			}, cm)).To(Succeed())
			Expect(cm.Data).To(HaveKey("metrics.json"))
			Expect(cm.Labels).To(HaveKeyWithValue("stats.corium.io/collector", collectorName))
		})
	})

	Context("When the referenced config does not exist", func() {
		const collectorName = "test-collector-no-config"
		ctx := context.Background()
		collectorNN := types.NamespacedName{Name: collectorName, Namespace: namespace}

		BeforeEach(func() {
			collector := &statsv1alpha1.JAXStatsCollector{}
			err := k8sClient.Get(ctx, collectorNN, collector)
			if err != nil && errors.IsNotFound(err) {
				collector = &statsv1alpha1.JAXStatsCollector{
					ObjectMeta: metav1.ObjectMeta{Name: collectorName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsCollectorSpec{
						TargetNamespace: namespace,
						ConfigRef:       "nonexistent-config",
						Metrics:         []string{"memory_usage"},
						Selector: metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "test"},
						},
					},
				}
				Expect(k8sClient.Create(ctx, collector)).To(Succeed())
			}
		})

		AfterEach(func() {
			collector := &statsv1alpha1.JAXStatsCollector{}
			if err := k8sClient.Get(ctx, collectorNN, collector); err == nil {
				k8sClient.Delete(ctx, collector)
			}
		})

		It("should set status to Error with message about missing config", func() {
			reconciler := &JAXStatsCollectorReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(10),
			}

			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: collectorNN})
			Expect(err).NotTo(HaveOccurred())

			updated := &statsv1alpha1.JAXStatsCollector{}
			Expect(k8sClient.Get(ctx, collectorNN, updated)).To(Succeed())
			Expect(updated.Status.CollectionStatus).To(Equal("Error"))
			Expect(updated.Status.ErrorMessage).To(ContainSubstring("not found"))
		})
	})

	Context("When selector matches no pods", func() {
		const (
			configName    = "collector-empty-config"
			collectorName = "test-collector-empty"
		)
		ctx := context.Background()
		collectorNN := types.NamespacedName{Name: collectorName, Namespace: namespace}

		BeforeEach(func() {
			config := &statsv1alpha1.JAXStatsConfig{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: configName, Namespace: namespace}, config)
			if err != nil && errors.IsNotFound(err) {
				config = &statsv1alpha1.JAXStatsConfig{
					ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsConfigSpec{
						Enabled:            true,
						CollectionInterval: 30,
						Metrics:            []string{"memory_usage"},
						StorageConfig:      statsv1alpha1.StorageConfig{Type: "configmap"},
					},
				}
				Expect(k8sClient.Create(ctx, config)).To(Succeed())
			}

			collector := &statsv1alpha1.JAXStatsCollector{}
			err = k8sClient.Get(ctx, collectorNN, collector)
			if err != nil && errors.IsNotFound(err) {
				collector = &statsv1alpha1.JAXStatsCollector{
					ObjectMeta: metav1.ObjectMeta{Name: collectorName, Namespace: namespace},
					Spec: statsv1alpha1.JAXStatsCollectorSpec{
						TargetNamespace: namespace,
						ConfigRef:       configName,
						Metrics:         []string{"memory_usage"},
						Selector: metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "nonexistent-app"},
						},
					},
				}
				Expect(k8sClient.Create(ctx, collector)).To(Succeed())
			}
		})

		AfterEach(func() {
			cm := &corev1.ConfigMap{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: collectorName + "-metrics", Namespace: namespace}, cm); err == nil {
				k8sClient.Delete(ctx, cm)
			}
			collector := &statsv1alpha1.JAXStatsCollector{}
			if err := k8sClient.Get(ctx, collectorNN, collector); err == nil {
				k8sClient.Delete(ctx, collector)
			}
			config := &statsv1alpha1.JAXStatsConfig{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: configName, Namespace: namespace}, config); err == nil {
				config.Finalizers = nil
				k8sClient.Update(ctx, config)
				k8sClient.Delete(ctx, config)
			}
		})

		It("should create ConfigMap with zero pods", func() {
			reconciler := &JAXStatsCollectorReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(10),
			}

			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: collectorNN})
			Expect(err).NotTo(HaveOccurred())

			updated := &statsv1alpha1.JAXStatsCollector{}
			Expect(k8sClient.Get(ctx, collectorNN, updated)).To(Succeed())
			Expect(updated.Status.CollectionStatus).To(Equal("Active"))
			Expect(updated.Status.DiscoveredPods).To(Equal(int32(0)))
		})
	})
})
