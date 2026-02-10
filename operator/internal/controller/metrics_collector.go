package controller

import (
	"encoding/json"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// PodMetrics holds collected metrics for a single pod.
type PodMetrics struct {
	Name           string `json:"name"`
	Namespace      string `json:"namespace"`
	Phase          string `json:"phase"`
	Ready          bool   `json:"ready"`
	RestartCount   int32  `json:"restart_count"`
	ContainerCount int    `json:"container_count"`
	StartTime      string `json:"start_time,omitempty"`
	NodeName       string `json:"node_name"`
	IP             string `json:"ip,omitempty"`
	CPURequest     string `json:"cpu_request,omitempty"`
	MemoryRequest  string `json:"memory_request,omitempty"`
}

// CollectedMetrics holds the full collection result persisted to ConfigMap.
type CollectedMetrics struct {
	CollectorName string         `json:"collector_name"`
	Namespace     string         `json:"namespace"`
	CollectedAt   string         `json:"collected_at"`
	PodCount      int            `json:"pod_count"`
	Pods          []PodMetrics   `json:"pods"`
	Summary       MetricsSummary `json:"summary"`
}

// MetricsSummary provides aggregate stats across all discovered pods.
type MetricsSummary struct {
	TotalPods     int   `json:"total_pods"`
	ReadyPods     int   `json:"ready_pods"`
	NotReadyPods  int   `json:"not_ready_pods"`
	TotalRestarts int32 `json:"total_restarts"`
}

// CollectPodMetrics is a pure function that extracts metrics from a list of pods.
func CollectPodMetrics(pods []corev1.Pod) []PodMetrics {
	result := make([]PodMetrics, 0, len(pods))
	for _, pod := range pods {
		pm := PodMetrics{
			Name:           pod.Name,
			Namespace:      pod.Namespace,
			Phase:          string(pod.Status.Phase),
			ContainerCount: len(pod.Spec.Containers),
			NodeName:       pod.Spec.NodeName,
			IP:             pod.Status.PodIP,
		}

		if pod.Status.StartTime != nil {
			pm.StartTime = pod.Status.StartTime.Format(time.RFC3339)
		}

		// Aggregate restart counts and readiness from container statuses
		var totalRestarts int32
		ready := true
		for _, cs := range pod.Status.ContainerStatuses {
			totalRestarts += cs.RestartCount
			if !cs.Ready {
				ready = false
			}
		}
		pm.RestartCount = totalRestarts
		pm.Ready = ready && pod.Status.Phase == corev1.PodRunning

		// Extract resource requests from the first container
		if len(pod.Spec.Containers) > 0 {
			reqs := pod.Spec.Containers[0].Resources.Requests
			if cpu, ok := reqs[corev1.ResourceCPU]; ok {
				pm.CPURequest = cpu.String()
			}
			if mem, ok := reqs[corev1.ResourceMemory]; ok {
				pm.MemoryRequest = mem.String()
			}
		}

		result = append(result, pm)
	}
	return result
}

// BuildCollectedMetrics creates the full metrics payload.
func BuildCollectedMetrics(collectorName, namespace string, podMetrics []PodMetrics) CollectedMetrics {
	summary := MetricsSummary{
		TotalPods: len(podMetrics),
	}
	for _, pm := range podMetrics {
		if pm.Ready {
			summary.ReadyPods++
		} else {
			summary.NotReadyPods++
		}
		summary.TotalRestarts += pm.RestartCount
	}

	return CollectedMetrics{
		CollectorName: collectorName,
		Namespace:     namespace,
		CollectedAt:   time.Now().UTC().Format(time.RFC3339),
		PodCount:      len(podMetrics),
		Pods:          podMetrics,
		Summary:       summary,
	}
}

// MarshalMetrics serializes the collected metrics to JSON.
func MarshalMetrics(cm CollectedMetrics) (string, error) {
	data, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
