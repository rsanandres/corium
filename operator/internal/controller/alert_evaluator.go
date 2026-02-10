package controller

import (
	"fmt"
	"strconv"

	statsv1alpha1 "github.com/raph/corium/operator/api/v1alpha1"
)

// AlertEvalResult holds the evaluation result for a single rule.
type AlertEvalResult struct {
	RuleName    string
	Firing      bool
	Severity    string
	Metric      string
	ActualValue float64
	Threshold   float64
	Message     string
}

// EvaluateAlertRules is a pure function that evaluates alert rules against collected metrics.
func EvaluateAlertRules(rules []statsv1alpha1.AlertRule, metrics *CollectedMetrics) []AlertEvalResult {
	if metrics == nil {
		return nil
	}

	results := make([]AlertEvalResult, 0, len(rules))
	for _, rule := range rules {
		result := evaluateRule(rule, metrics)
		results = append(results, result)
	}
	return results
}

func evaluateRule(rule statsv1alpha1.AlertRule, metrics *CollectedMetrics) AlertEvalResult {
	threshold, err := strconv.ParseFloat(rule.Threshold, 64)
	if err != nil {
		return AlertEvalResult{
			RuleName: rule.Name,
			Firing:   false,
			Severity: rule.Severity,
			Metric:   rule.Metric,
			Message:  fmt.Sprintf("invalid threshold value: %s", rule.Threshold),
		}
	}

	actual := getMetricValue(rule.Metric, metrics)

	firing := compareValues(actual, threshold, rule.Operator)

	msg := fmt.Sprintf("%s: %s %.0f %s %.0f", rule.Name, rule.Metric, actual, rule.Operator, threshold)
	if firing {
		msg = fmt.Sprintf("FIRING: %s (actual=%.0f, threshold=%.0f)", rule.Name, actual, threshold)
	}

	return AlertEvalResult{
		RuleName:    rule.Name,
		Firing:      firing,
		Severity:    rule.Severity,
		Metric:      rule.Metric,
		ActualValue: actual,
		Threshold:   threshold,
		Message:     msg,
	}
}

func getMetricValue(metric string, metrics *CollectedMetrics) float64 {
	switch metric {
	case "restart_count":
		return float64(metrics.Summary.TotalRestarts)
	case "not_ready_count":
		return float64(metrics.Summary.NotReadyPods)
	case "container_count":
		var total int
		for _, p := range metrics.Pods {
			total += p.ContainerCount
		}
		return float64(total)
	case "pod_count":
		return float64(metrics.Summary.TotalPods)
	default:
		return 0
	}
}

func compareValues(actual, threshold float64, op string) bool {
	switch op {
	case ">":
		return actual > threshold
	case "<":
		return actual < threshold
	case ">=":
		return actual >= threshold
	case "<=":
		return actual <= threshold
	case "==":
		return actual == threshold
	default:
		return false
	}
}
