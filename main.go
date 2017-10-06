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

var (
	dockerHost                    string
	sweepInterval                 int
	removeImages                  bool
	removeHealthyContainersExited bool
)

func init() {
	flag.StringVar(&dockerHost, "docker-host", "unix:///var/run/docker.sock", "-docker-host=unix:///var/run/docker.sock")
	flag.IntVar(&sweepInterval, "sweep-interval", 60, "-sweep-interval=60")
	flag.BoolVar(&removeImages, "remove-images", true, "-remove-images=true")
	flag.BoolVar(&removeHealthyContainersExited, "remove-healthy-containers-exited", true, "-remove-healthy-containers-exited=true")
}

func main() {
	flag.Parse()

	dc, err := docker.NewClient(dockerHost)
	if err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGTERM, syscall.SIGINT)

	fmt.Fprintf(os.Stdout, "gcd: [info]: (Time: %v)\n", time.Now().String())
	fmt.Fprintf(os.Stdout, "gcd: [info]: (Docker Host: %v)\n", dockerHost)
	fmt.Fprintf(os.Stdout, "gcd: [info]: (Sweep Interval: %vs)\n", sweepInterval)
	fmt.Fprintf(os.Stdout, "gcd: [info]: (Remove Images: %v)\n", removeImages)
	fmt.Fprintf(os.Stdout, "gcd: [info]: (Remove Healthy Containers Exited: %v)\n", removeHealthyContainersExited)

	for {
		select {
		case <-s:
			os.Exit(0)
		case <-time.Tick(time.Duration(sweepInterval) * time.Second):
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
						if (removeHealthyContainersExited && exitCodeFromContainer == "(0)") || exitCodeFromContainer != "(0)" {
							fmt.Fprintf(os.Stdout, "gcd: [removing container]: (Id: %v, Labels: %v)\n", container.ID, container.Labels)
							if err := dc.RemoveContainer(docker.RemoveContainerOptions{
								ID:            container.ID,
								RemoveVolumes: true,
								Force:         true,
							}); err == nil {
								fmt.Fprintf(os.Stdout, "gcd: [removed container]: (Id: %v, Labels: %v)\n", container.ID, container.Labels)
							}
						}
					}
				}()
			}
			wgContainers.Wait()
			if removeImages {
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
						if err := dc.RemoveImage(image.ID); err == nil {
							fmt.Fprintf(os.Stdout, "gcd: [removed image]: (Id: %v, Tags: %v)\n", image.ID, image.RepoTags)
						}
					}()
				}
				wgImages.Wait()
			}
		}
	}
}
