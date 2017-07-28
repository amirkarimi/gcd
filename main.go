package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/stone-payments/gcd/worker"
)

func main() {
	interval := 1

	host := os.Getenv("GCD_DOCKER_HOST")
	if host == "" {
		host = "unix:///var/run/docker.sock"
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

	worker, err := worker.New(host, "1.24")
	if err != nil {
		panic(err)
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-sig:
			fmt.Printf("\nDown daemon by signal: %v\n", <-sig)
			os.Exit(0)
		case <-time.After(time.Second * time.Duration(interval)):
			worker.Sweep()
		}
	}
}
