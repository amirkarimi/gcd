package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
)

var (
	target                 string
	sweepInterval          int
	removeImages           bool
	removeHealthyContainer bool
)

func init() {
	flag.StringVar(&target, "target", "unix:///var/run/docker.sock", "-target=unix:///var/run/docker.sock")
	flag.IntVar(&sweepInterval, "sweep-interval", 1, "-sweep-interval=1")
	flag.BoolVar(&removeImages, "remove-images", true, "-remove-images=true")
	flag.BoolVar(&removeHealthyContainer, "remove-healthy-container", true, "-remove-healthy-container=true")
}

func main() {
	flag.Parse()

	dc, err := docker.NewClient(target)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.Out = os.Stdout

	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)

	logger.Infof("Target: %v", target)
	logger.Infof("Sweep Interval: %vs", sweepInterval)
	logger.Infof("Remove Images: %v", removeImages)
	logger.Infof("Remove Healthy Containers: %v", removeHealthyContainer)

	for {
		select {
		case <-s:
			logger.Fatalf("Down worker by signal: %v", <-s)
		case <-time.Tick(time.Duration(sweepInterval) * time.Second):
			logger.Infof("Time: %v", time.Now().UnixNano())
			containers, err := dc.ListContainers(docker.ListContainersOptions{
				All: true,
			})
			if err != nil {
				logger.Error(err)
			}
			for _, container := range containers {
				exitCodeFromContainer := "(-)"
				if splitedStatus := strings.Split(container.Status, " "); len(splitedStatus) > 1 {
					exitCodeFromContainer = splitedStatus[1]
				}
				if removeHealthyContainer && exitCodeFromContainer == "(0)" && container.State == "exited" {
					err := dc.RemoveContainer(docker.RemoveContainerOptions{
						ID:            container.ID,
						RemoveVolumes: true,
						Force:         true,
					})
					if err != nil {
						logger.Errorf("[Remove Container]: %v", err)
					} else {
						logger.Infof("[Remove Container]: %v", container.ID)
					}
				}
			}

			if removeImages {
				images, err := dc.ListImages(docker.ListImagesOptions{})
				if err != nil {
					logger.Error(err)
				}
				for _, image := range images {
					err := dc.RemoveImage(image.ID)
					if err != nil {
						logger.Errorf("[Remove Image]: %v", err)
					} else {
						logger.Infof("[Remove Image]: %v", image.ID)
					}
				}
			}
		}
		fmt.Println()
	}
}
