package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	promNamespace = "radeon_exporter"

	state = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "state",
		Help:      "normal / fault",
	}, []string{"model_name"})

	PowerSuppliedBy = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "power_supplied_by",
		Help:      "utility / battery",
	}, []string{"model_name"})

	UtilityVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "utility_voltage",
		Help:      "Utility Voltage",
	}, []string{"model_name"})

	OutputVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "output_voltage",
		Help:      "Output Voltage",
	}, []string{"model_name"})

	BatteryCapacity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "battery_capacity",
		Help:      "Battery Capacity as %",
	}, []string{"model_name"})

	RemainingRuntime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "remaining_runtime",
		Help:      "Remaining Runtime on battery",
	}, []string{"model_name"})

	LoadWatts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "load_watts",
		Help:      "Current Load in watts",
	}, []string{"model_name"})

	LoadPct = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "load_pct",
		Help:      "current load as %",
	}, []string{"model_name"})

	LineInteraction = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "line_interaction",
		Help:      "ups line interaction",
	}, []string{"model_name"})

	TestResult = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "test_result",
		Help:      "result of last test result",
	}, []string{"model_name"})

	TestResultTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "test_result_time",
		Help:      "when the last test was run",
	}, []string{"model_name"})

	LastPowerEvent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "last_power_event",
		Help:      "why power went out",
	}, []string{"model_name"})

	LastPowerEventTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "last_power_event_time",
		Help:      "the last time power went out",
	}, []string{"model_name"})

	LastPowerEventDuration = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Name:      "last_power_event_duration",
		Help:      "how long the last event lasted",
	}, []string{"model_name"})
)
