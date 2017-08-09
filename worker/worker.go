package worker

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/stone-payments/gcd/clients/docker"
	"github.com/stone-payments/gcd/logger"
)

type (
	Worker interface {
		ListContainers(ch chan []docker.Container)
		ListImages(ch chan []docker.Image)
		RemoveContainer(wg *sync.WaitGroup, container docker.Container)
		RemoveImage(wg *sync.WaitGroup, image docker.Image)
		Sweep()
		GetVersion() string
	}

	WorkerAsClient struct {
		dockerClient           docker.Client
		logger                 logger.Logger
		removeImages           bool
		removeContainersExited bool
	}
)

func (wac WorkerAsClient) ListContainers(ch chan []docker.Container) {
	containers, err := wac.dockerClient.GetContainers()
	if err != nil {
		wac.logger.Error(err.Error())
	}

	ch <- containers
}

func (wac WorkerAsClient) ListImages(ch chan []docker.Image) {
	images, err := wac.dockerClient.GetImages()
	if err != nil {
		wac.logger.Error(err.Error())
	}

	ch <- images
}

func (wac WorkerAsClient) RemoveContainer(wg *sync.WaitGroup, container docker.Container) {
	if container.State != "running" {
		isNormalExited := strings.Contains(container.Status, "Exited (0)")
		if (isNormalExited && wac.removeContainersExited) || !isNormalExited {
			ok, err := wac.dockerClient.RemoveContainer(container.Id)
			if err != nil {
				wac.logger.Error(err.Error())
			} else {
				if ok {
					wac.logger.OK("Container", container.Id, "removed successful")
				} else {
					wac.logger.Skip("Container", container.Id)
				}
			}
		} else {
			wac.logger.Skip("Container", container.Id, "skipped, Status:", container.Status)
		}
	}
	wg.Done()
}

func (wac WorkerAsClient) RemoveImage(wg *sync.WaitGroup, image docker.Image) {
	if ok, err := wac.dockerClient.RemoveImage(image.Id); err != nil {
		wac.logger.Error(err.Error())
	} else {
		if ok {
			wac.logger.OK("Image", image.Id, "removed successful")
		} else {
			wac.logger.Skip("Image", image.Id)
		}
	}
	wg.Done()
}

func (wac WorkerAsClient) Sweep() {
	fmt.Println("")

	chContainers := make(chan []docker.Container)
	chImages := make(chan []docker.Image)

	go wac.ListContainers(chContainers)
	go wac.ListImages(chImages)

	select {
	case containers := <-chContainers:
		wac.logger.Info("Time:", time.Now().UnixNano())
		wac.logger.Info("Host:", wac.dockerClient.GetHost())
		wac.logger.Info("Containers total:", len(containers))
		select {
		case images := <-chImages:
			wac.logger.Info("Images total:", len(images))

			wgContainerDone := make(chan bool)
			var wgContainer sync.WaitGroup
			for _, container := range containers {
				wgContainer.Add(1)
				go wac.RemoveContainer(&wgContainer, container)
			}

			go func(done chan bool) {
				wgContainer.Wait()
				done <- true
			}(wgContainerDone)

			select {
			case <-wgContainerDone:
				wac.logger.Info("Action to containers finished")
			case <-time.After(time.Second * 1):
				wac.logger.Error("Action to containers timeout occurred")
			}

			if wac.removeImages {
				wgImageDone := make(chan bool)
				var wgImage sync.WaitGroup
				for _, image := range images {
					wgImage.Add(1)
					go wac.RemoveImage(&wgImage, image)
				}

				go func(done chan bool) {
					wgImage.Wait()
					done <- true
				}(wgImageDone)

				select {
				case <-wgImageDone:
					wac.logger.Info("Action to images finished")
				case <-time.After(time.Second * 1):
					wac.logger.Error("Action to images timeout occurred")
				}
			} else {
				wac.logger.Info("Remove images:", wac.removeImages)
			}

		case <-time.After(time.Second * 1):
			wac.logger.Error("Images request timeout")
		}

	case <-time.After(time.Second * 1):
		wac.logger.Error("Containers request timeout")
	}
}

func (wac WorkerAsClient) GetVersion() string {
	return wac.dockerClient.GetVersion()
}

func New(dockerClient docker.Client, logger logger.Logger, removeImage, removeContainersPaused bool) Worker {
	return WorkerAsClient{
		dockerClient,
		logger,
		removeImage,
		removeContainersPaused,
	}
}
