package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const modelLabel = "model_name"

// nolint: gochecknoglobals
var (
	promNamespace = "cyber_power_exporter"

	stateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "state",
		Help:      "0=Normal / 1=Power Failure",
	}, []string{modelLabel})

	powerSuppliedByGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "power_supplied_by",
		Help:      "0=Utility Power / 1=Battery Power",
	}, []string{modelLabel})

	utilityVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "utility_voltage",
		Help:      "Utility Voltage",
	}, []string{modelLabel})

	outputVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "output_voltage",
		Help:      "Output Voltage",
	}, []string{modelLabel})

	batteryCapacityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "battery_capacity",
		Help:      "Battery Capacity as %",
	}, []string{modelLabel})

	remainingRuntimeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "remaining_runtime",
		Help:      "Remaining Runtime on battery in seconds",
	}, []string{modelLabel})

	loadWattsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "load_watts",
		Help:      "Current Load in watts",
	}, []string{modelLabel})

	loadPctGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "load_pct",
		Help:      "current load as %",
	}, []string{modelLabel})

	lineInteractionGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "line_interaction",
		Help:      "ups line interaction",
	}, []string{modelLabel})

	testResultGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "test_result",
		Help:      "result of last test result",
	}, []string{modelLabel})

	lastPowerEventDurationGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "last_power_event_duration",
		Help:      "how long the last event lasted",
	}, []string{modelLabel})
)
