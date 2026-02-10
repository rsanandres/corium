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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	statsv1alpha1 "github.com/raph/corium/operator/api/v1alpha1"
)

var _ = Describe("JAXStatsConfig Controller", func() {
	const namespace = "default"

	Context("When reconciling a valid enabled config", func() {
		const resourceName = "test-config-valid"
		ctx := context.Background()
		namespacedName := types.NamespacedName{Name: resourceName, Namespace: namespace}

		BeforeEach(func() {
			resource := &statsv1alpha1.JAXStatsConfig{}
			err := k8sClient.Get(ctx, namespacedName, resource)
			if err != nil && errors.IsNotFound(err) {
				resource = &statsv1alpha1.JAXStatsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: namespace,
					},
					Spec: statsv1alpha1.JAXStatsConfigSpec{
						Enabled:            true,
						CollectionInterval: 30,
						Metrics:            []string{"memory_usage"},
						StorageConfig: statsv1alpha1.StorageConfig{
							Type:     "prometheus",
							Endpoint: "http://prometheus:9090",
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &statsv1alpha1.JAXStatsConfig{}
			err := k8sClient.Get(ctx, namespacedName, resource)
			if err == nil {
				resource.Finalizers = nil
				Expect(k8sClient.Update(ctx, resource)).To(Succeed())
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should set status to Active and add finalizer", func() {
			reconciler := &JAXStatsConfigReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(10),
			}

			result, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeNumerically(">", 0))

			// Reconcile again (finalizer was just added, now real logic runs)
			result, err = reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			Expect(err).NotTo(HaveOccurred())

			updated := &statsv1alpha1.JAXStatsConfig{}
			Expect(k8sClient.Get(ctx, namespacedName, updated)).To(Succeed())
			Expect(updated.Status.CollectionStatus).To(Equal("Active"))
			Expect(updated.Finalizers).To(ContainElement(configFinalizerName))
		})
	})

	Context("When reconciling a disabled config", func() {
		const resourceName = "test-config-disabled"
		ctx := context.Background()
		namespacedName := types.NamespacedName{Name: resourceName, Namespace: namespace}

		BeforeEach(func() {
			resource := &statsv1alpha1.JAXStatsConfig{}
			err := k8sClient.Get(ctx, namespacedName, resource)
			if err != nil && errors.IsNotFound(err) {
				resource = &statsv1alpha1.JAXStatsConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: namespace,
					},
					Spec: statsv1alpha1.JAXStatsConfigSpec{
						Enabled:            false,
						CollectionInterval: 60,
						StorageConfig: statsv1alpha1.StorageConfig{
							Type: "configmap",
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &statsv1alpha1.JAXStatsConfig{}
			err := k8sClient.Get(ctx, namespacedName, resource)
			if err == nil {
				resource.Finalizers = nil
				Expect(k8sClient.Update(ctx, resource)).To(Succeed())
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("should set status to Disabled", func() {
			reconciler := &JAXStatsConfigReconciler{
				Client:   k8sClient,
				Scheme:   k8sClient.Scheme(),
				Recorder: record.NewFakeRecorder(10),
			}

			// First reconcile adds finalizer
			reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			// Second reconcile runs main logic
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: namespacedName})
			Expect(err).NotTo(HaveOccurred())

			updated := &statsv1alpha1.JAXStatsConfig{}
			Expect(k8sClient.Get(ctx, namespacedName, updated)).To(Succeed())
			Expect(updated.Status.CollectionStatus).To(Equal("Disabled"))
		})
	})

	Context("When creating a config with invalid storage type", func() {
		It("should be rejected by API validation", func() {
			resource := &statsv1alpha1.JAXStatsConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-config-invalid-storage",
					Namespace: namespace,
				},
				Spec: statsv1alpha1.JAXStatsConfigSpec{
					Enabled:            true,
					CollectionInterval: 60,
					StorageConfig:      statsv1alpha1.StorageConfig{},
				},
			}
			err := k8sClient.Create(context.Background(), resource)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsInvalid(err)).To(BeTrue())
		})
	})
})
