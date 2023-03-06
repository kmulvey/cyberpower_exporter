package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	var cmdPath string
	var pollInterval time.Duration
	var v, db, http, prom bool
	flag.StringVar(&cmdPath, "cmd-path", "/usr/sbin/pwrstat", "absolute path to pwstat command")
	flag.DurationVar(&pollInterval, "poll-interval", time.Minute, "time interval to gather power stats")
	flag.BoolVar(&db, "db", false, "print version")
	flag.BoolVar(&http, "http", false, "print version")
	flag.BoolVar(&prom, "prom", false, "print version")
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
		dbHandle, err = dbConnect("cp.db")
		if err != nil {
			log.Fatal(err)
		}
	}

	if http {
		go webServer(":8080", dbHandle)
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

			status, _, err := parsePowerStats(out) // DEVICE
			if err != nil {
				log.Fatal(err)
			}

			if db {
				if err := insert(dbHandle, status); err != nil {
					log.Fatal(err)
				}
			}

			if prom {

			}

		case <-sigChannel:
			log.Info("shutting down")
			return
		}
	}
}
