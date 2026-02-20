package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	discoveredPodsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "corium_discovered_pods",
			Help: "Number of pods discovered by each CoriumMonitorCollector",
		},
		[]string{"collector"},
	)

	activeAlertsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "corium_active_alerts",
			Help: "Number of currently firing alerts per CoriumMonitorAlert",
		},
		[]string{"alert"},
	)

	reconcileErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "corium_reconcile_errors_total",
			Help: "Total number of reconciliation errors by controller",
		},
		[]string{"controller"},
	)
)

func init() {
	metrics.Registry.MustRegister(discoveredPodsGauge, activeAlertsGauge, reconcileErrorsTotal)
}
