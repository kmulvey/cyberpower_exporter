package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"go.szostok.io/version"
	"go.szostok.io/version/printer"
)

// nolint: gochecknoglobals
var dateFormat = "2006/01/02 15:04:05"

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: dateFormat,
	})

	var sigChannel = make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	// get user opts
	var cmdPath, promAddr string
	var pollInterval time.Duration
	var v bool
	flag.StringVar(&cmdPath, "cmd-path", "/usr/sbin/pwrstat", "absolute path to pwstat command")
	flag.StringVar(&promAddr, "prom-addr", ":9300", "bind address of the prom http server")
	flag.DurationVar(&pollInterval, "poll-interval", time.Second*5, "time interval to gather power stats")
	flag.BoolVar(&v, "version", false, "print version")
	flag.BoolVar(&v, "v", false, "print version")

	flag.Parse()

	if v {
		var verPrinter = printer.New()
		var info = version.Get()
		if err := verPrinter.PrintInfo(os.Stdout, info); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())

		var server = &http.Server{
			Addr:         promAddr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		if err := server.ListenAndServe(); err != nil {
			log.Fatal("http server error: ", err)
		}
	}()
	log.Info("started, go to grafana to monitor")

	gatherAndSaveStats(cmdPath)

	var ticker = time.NewTicker(pollInterval)
	for {
		select {
		case <-ticker.C:
			gatherAndSaveStats(cmdPath)

		case <-sigChannel:
			log.Info("shutting down")
			return
		}
	}
}

func gatherAndSaveStats(cmdPath string) {
	out, err := getPowerStats(cmdPath)
	if err != nil {
		log.Error(err)
	}

	status, err := parsePowerStatus(out)
	if err != nil {
		log.Error(err)
	}

	device, err := parseDeviceProperties(out)
	if err != nil {
		log.Error(err)
	}

	switch status.State {
	case "Normal":
		stateGauge.WithLabelValues(device.ModelName).Set(0)
	case "Power Failure":
		stateGauge.WithLabelValues(device.ModelName).Set(1)
	}

	switch status.PowerSupplyBy {
	case "Utility Power":
		powerSuppliedByGauge.WithLabelValues(device.ModelName).Set(0)
	case "Battery Power":
		powerSuppliedByGauge.WithLabelValues(device.ModelName).Set(1)
	}

	if status.LineInteraction == "None" {
		lineInteractionGauge.WithLabelValues(device.ModelName).Set(0)
	} else {
		lineInteractionGauge.WithLabelValues(device.ModelName).Set(1)
	}

	if status.TestResult == "Passed" {
		testResultGauge.WithLabelValues(device.ModelName).Set(0)
	} else {
		testResultGauge.WithLabelValues(device.ModelName).Set(1)
	}

	utilityVoltageGauge.WithLabelValues(device.ModelName).Set(float64(status.UtilityVoltage))
	outputVoltageGauge.WithLabelValues(device.ModelName).Set(float64(status.OutputVoltage))
	batteryCapacityGauge.WithLabelValues(device.ModelName).Set(float64(status.BatteryCapacity))
	remainingRuntimeGauge.WithLabelValues(device.ModelName).Set(float64(status.RemainingRuntime.Seconds()))
	loadWattsGauge.WithLabelValues(device.ModelName).Set(float64(status.LoadWatts))
	loadPctGauge.WithLabelValues(device.ModelName).Set(float64(status.LoadPct))
	lastPowerEventDurationGauge.WithLabelValues(device.ModelName).Set(float64(status.LastPowerEventDuration.Seconds()))
}
