package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	discoveredPodsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jaxstats_discovered_pods",
			Help: "Number of pods discovered by each JAXStatsCollector",
		},
		[]string{"collector"},
	)

	activeAlertsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jaxstats_active_alerts",
			Help: "Number of currently firing alerts per JAXStatsAlert",
		},
		[]string{"alert"},
	)

	reconcileErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "jaxstats_reconcile_errors_total",
			Help: "Total number of reconciliation errors by controller",
		},
		[]string{"controller"},
	)
)

func init() {
	metrics.Registry.MustRegister(discoveredPodsGauge, activeAlertsGauge, reconcileErrorsTotal)
}
