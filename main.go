package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

type Options struct {
	dockerHost                    string
	sweepInterval                 int
	removeImages                  bool
	removeHealthyContainersExited bool
}

var options Options

func init() {
	flag.StringVar(&options.dockerHost, "docker-host", "unix:///var/run/docker.sock", "-docker-host=unix:///var/run/docker.sock")
	flag.IntVar(&options.sweepInterval, "sweep-interval", 60, "-sweep-interval=60")
	flag.BoolVar(&options.removeImages, "remove-images", true, "-remove-images=true")
	flag.BoolVar(&options.removeHealthyContainersExited, "remove-healthy-containers-exited", true, "-remove-healthy-containers-exited=true")
}

func main() {
	flag.Parse()

	dc, err := docker.NewClient(options.dockerHost)
	if err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)

	fmt.Fprintf(os.Stdout, "gcd: [info]: (Time: %v)\n", time.Now().String())
	fmt.Fprintf(os.Stdout, "gcd: [info]: (Docker Host: %v)\n", options.dockerHost)
	fmt.Fprintf(os.Stdout, "gcd: [info]: (Sweep Interval: %vs)\n", options.sweepInterval)
	fmt.Fprintf(os.Stdout, "gcd: [info]: (Remove Images: %v)\n", options.removeImages)
	fmt.Fprintf(os.Stdout, "gcd: [info]: (Remove Healthy Containers Exited: %v)\n", options.removeHealthyContainersExited)

	for {
		select {
		case <-s:
			os.Exit(0)
		case <-time.Tick(time.Duration(options.sweepInterval) * time.Second):
			Run(dc, options)
		}
	}
}

func Run(dc *docker.Client, options Options) {
	fmt.Fprintf(os.Stdout, "\ngcd: [info]: (Time: %v)\n", time.Now().String())
	containers, err := dc.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	var wgContainers sync.WaitGroup
	for _, container := range containers {
		wgContainers.Add(1)
		go func() {
			defer wgContainers.Done()
			exitCodeFromContainer := "(-)"
			if splitedStatus := strings.Split(container.Status, " "); len(splitedStatus) > 1 {
				exitCodeFromContainer = splitedStatus[1]
			}
			if container.State != "running" {
				if (options.removeHealthyContainersExited && exitCodeFromContainer == "(0)") || exitCodeFromContainer != "(0)" {
					fmt.Fprintf(os.Stdout, "gcd: [trying remove container]: (Id: %v, Labels: %v)\n", container.ID, container.Labels)
					if err := dc.RemoveContainer(docker.RemoveContainerOptions{
						ID:            container.ID,
						RemoveVolumes: true,
						Force:         true,
					}); err != nil {
						fmt.Fprintf(os.Stderr, "gcd: [error]: (when try remove container, reason: %v)\n", err.Error())
					} else {
						fmt.Fprintf(os.Stdout, "gcd: [removed container]: (Id: %v, Labels: %v)\n", container.ID, container.Labels)
					}
				}
			}
		}()
	}
	wgContainers.Wait()
	if options.removeImages {
		var wgImages sync.WaitGroup
		images, err := dc.ListImages(docker.ListImagesOptions{})
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
		for _, image := range images {
			wgImages.Add(1)
			go func() {
				defer wgImages.Done()
				fmt.Fprintf(os.Stdout, "gcd: [trying remove image]: (Id: %v, Tags: %v)\n", image.ID, image.RepoTags)
				if err := dc.RemoveImage(image.ID); err != nil {
					fmt.Fprintf(os.Stderr, "gcd: [error]: (when try remove image, reason: %v)\n", err.Error())
				} else {
					fmt.Fprintf(os.Stdout, "gcd: [removed image]: (Id: %v, Tags: %v)\n", image.ID, image.RepoTags)
				}
			}()
		}
		wgImages.Wait()
	}
}
