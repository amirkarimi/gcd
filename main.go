package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/stone-payments/gcd/clients/docker"
	"github.com/stone-payments/gcd/logger"
	"github.com/stone-payments/gcd/worker"
)

func main() {
	interval := 1
	removeImages := true
	removeContainersExited := false // Enable remove to containers with exited status 0

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

	if removeImagesString := os.Getenv("GCD_REMOVE_IMAGES"); removeImagesString != "" {
		value, err := strconv.ParseBool(removeImagesString)
		if err != nil {
			fmt.Println("Invalid value as option for remove image:", err.Error())
		} else {
			removeImages = value
		}
	}

	if removeContainersExitedString := os.Getenv("GCD_REMOVE_CONTAINERS_EXITED"); removeContainersExitedString != "" {
		value, err := strconv.ParseBool(removeContainersExitedString)
		if err != nil {
			fmt.Println("Invalid value as option for remove containers paused:", err.Error())
		} else {
			removeContainersExited = value
		}
	}

	logger := logger.New()

	dockerClient, err := docker.NewClient(host, version)
	if err != nil {
		logger.Exit(1, err.Error())
	}

	w := worker.New(dockerClient, logger, removeImages, removeContainersExited)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	logger.Info("Time:", time.Now().UnixNano())
	logger.Info("State:", "running")
	logger.Info("Docker API Version:", dockerClient.GetVersion())

	for {
		select {
		case <-sig:
			logger.Exit(0, "Down daemon by signal:", <-sig)
		case <-time.After(time.Second * time.Duration(interval)):
			w.Sweep()
		}
	}
}
