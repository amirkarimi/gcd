package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/stone-payments/gcd/logger"
	"github.com/stone-payments/gcd/worker"
)

func main() {
	interval := 1

	host := os.Getenv("GCD_DOCKER_HOST")
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}

	version := os.Getenv("GCD_DOCKER_API_VERSION")
	if version == "" {
		version = "1.24"
	}

	if intervalString := os.Getenv("GCD_SWEEP_INTERVAL"); intervalString != "" {
		value, err := strconv.Atoi(intervalString)
		if err != nil {
			fmt.Println("Invalid value as interval:", err.Error())
			interval = 1
		} else {
			interval = value
		}
	}

	logger := logger.New()

	worker, err := worker.New(host, version, logger)
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	logger.Info("Time:", time.Now().UnixNano())
	logger.Info("State:", "running")
	logger.Info("Docker API Version:", worker.GetVersion())

	for {
		select {
		case <-sig:
			logger.Exit("Down daemon by signal:", <-sig)
		case <-time.After(time.Second * time.Duration(interval)):
			worker.Sweep()
		}
	}
}
