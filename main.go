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
	"gorm.io/gorm"

	"go.szostok.io/version"
	"go.szostok.io/version/printer"
)

var dateFormat = "2006/01/02 15:04:05"

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: dateFormat,
	})

	var sigChannel = make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	// get user opts
	var cmdPath, httpAddr, promAddr, dbName string
	var pollInterval time.Duration
	var v, db, enableHttp, enableProm bool
	flag.StringVar(&cmdPath, "cmd-path", "/usr/sbin/pwrstat", "absolute path to pwstat command")
	flag.StringVar(&httpAddr, "http-addr", ":1000", "bind address of the http server")
	flag.StringVar(&cmdPath, "prom-addr", ":1001", "bind address of the prom http server")
	flag.StringVar(&dbName, "db-name", "cp.db", "name of the sqlite file")
	flag.DurationVar(&pollInterval, "poll-interval", time.Minute, "time interval to gather power stats")
	flag.BoolVar(&db, "db", false, "write to sqlite db")
	flag.BoolVar(&enableHttp, "http", false, "turn of http server")
	flag.BoolVar(&enableProm, "prom", false, "enable prom stats")
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

	var dbHandle *gorm.DB
	var err error

	if db {
		dbHandle, err = dbConnect(dbName)
		if err != nil {
			log.Fatal(err)
		}
	}

	if enableHttp {
		go webServer(httpAddr, dbHandle)
	}

	if enableProm {
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
	}

	var ticker = time.NewTicker(pollInterval)
	for {
		select {
		case <-ticker.C:
			log.Info("gathering stats")

			out, err := getPowerStats(cmdPath)
			if err != nil {
				log.Fatal(err)
			}

			status, device, err := parsePowerStats(out)
			if err != nil {
				log.Fatal(err)
			}

			if db {
				if err := insert(dbHandle, status); err != nil {
					log.Fatal(err)
				}
			}

			if enableProm {
				if status.State == "Normal" {
					stateGauge.WithLabelValues(device.ModelName).Set(0)
				} else if status.State == "Power Failure" {
					stateGauge.WithLabelValues(device.ModelName).Set(1)
				}

				if status.PowerSupplyBy == "Utility Power" {
					powerSuppliedByGauge.WithLabelValues(device.ModelName).Set(0)
				} else if status.State == "Battery Power" {
					powerSuppliedByGauge.WithLabelValues(device.ModelName).Set(1)
				}

				utilityVoltageGauge.WithLabelValues(device.ModelName).Set(float64(status.UtilityVoltage))
				outputVoltageGauge.WithLabelValues(device.ModelName).Set(float64(status.OutputVoltage))
				batteryCapacityGauge.WithLabelValues(device.ModelName).Set(float64(status.BatteryCapacity))
				remainingRuntimeGauge.WithLabelValues(device.ModelName).Set(float64(status.RemainingRuntime.Seconds()))
				loadWattsGauge.WithLabelValues(device.ModelName).Set(float64(status.LoadWatts))
				loadPctGauge.WithLabelValues(device.ModelName).Set(float64(status.LoadPct))
			}

		case <-sigChannel:
			log.Info("shutting down")
			return
		}
	}
}
