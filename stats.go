package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promNamespace = "cyber_power_exporter"

	stateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "state",
		Help:      "0=Normal / 1=Power Failure",
	}, []string{"model_name"})

	powerSuppliedByGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "power_supplied_by",
		Help:      "0=Utility Power / 1=Battery Power",
	}, []string{"model_name"})

	utilityVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "utility_voltage",
		Help:      "Utility Voltage",
	}, []string{"model_name"})

	outputVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "output_voltage",
		Help:      "Output Voltage",
	}, []string{"model_name"})

	batteryCapacityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "battery_capacity",
		Help:      "Battery Capacity as %",
	}, []string{"model_name"})

	remainingRuntimeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "remaining_runtime",
		Help:      "Remaining Runtime on battery in seconds",
	}, []string{"model_name"})

	loadWattsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "load_watts",
		Help:      "Current Load in watts",
	}, []string{"model_name"})

	loadPctGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "load_pct",
		Help:      "current load as %",
	}, []string{"model_name"})

	lineInteractionGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "line_interaction",
		Help:      "ups line interaction",
	}, []string{"model_name"})

	testResultGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "test_result",
		Help:      "result of last test result",
	}, []string{"model_name"})

	/*
		testResultTimeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: promNamespace,
			Name:      "test_result_time",
			Help:      "when the last test was run",
		}, []string{"model_name"})

		lastPowerEventGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: promNamespace,
			Name:      "last_power_event",
			Help:      "why power went out",
		}, []string{"model_name"})

		lastPowerEventTimeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: promNamespace,
			Name:      "last_power_event_time",
			Help:      "the last time power went out",
		}, []string{"model_name"})
	*/
	lastPowerEventDurationGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "last_power_event_duration",
		Help:      "how long the last event lasted",
	}, []string{"model_name"})
)
