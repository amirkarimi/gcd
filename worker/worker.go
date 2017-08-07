package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/stone-payments/gcd/logger"
)

type (
	Container struct {
		Id      string
		Created int
		SizeRw  int
		Status  string
		State   string
	}

	Image struct {
		Id string
	}

	Worker interface {
		ListContainers(ch chan []Container)
		ListImages(ch chan []Image)
		RemoveContainer(wg *sync.WaitGroup, container Container)
		RemoveImage(wg *sync.WaitGroup, image Image)
		Sweep()
		GetVersion() string
	}

	WorkerAsClient struct {
		conn                   http.Client
		host                   string
		version                string
		logger                 logger.Logger
		removeImages           bool
		removeContainersExited bool
	}
)

func (wac WorkerAsClient) ListContainers(ch chan []Container) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://v%v/containers/json?all=1", wac.version), nil)
	if err != nil {
		panic(err)
	}

	resp, err := wac.conn.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	containers := make([]Container, 0)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &containers); err != nil {
		panic(err)
	}

	ch <- containers
}

func (wac WorkerAsClient) ListImages(ch chan []Image) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://v%v/images/json?all=1", wac.version), nil)
	if err != nil {
		panic(err)
	}

	resp, err := wac.conn.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("", err.Error())
		panic(err)
	}

	images := make([]Image, 0)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &images); err != nil {
		panic(err)
	}

	ch <- images
}

func (wac WorkerAsClient) RemoveContainer(wg *sync.WaitGroup, container Container) {
	isNormalExited := strings.Contains(container.Status, "Exited (0)")
	if container.State != "running" {
		if (isNormalExited && wac.removeContainersExited) || !isNormalExited {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("http://v%v/containers/%v", wac.version, container.Id), nil)
			if err != nil {
				panic(err)
			}

			resp, err := wac.conn.Do(req)
			if err != nil {
				fmt.Sprintln(err.Error())
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNoContent {
				wac.logger.OK("Container", container.Id, "removed successful")
			}
		} else {
			wac.logger.Skip("Container", container.Id, "skipped, Status:", container.Status)
		}
	}
	wg.Done()
}

func (wac WorkerAsClient) RemoveImage(wg *sync.WaitGroup, image Image) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://v%v/images/%v?force=true", wac.version, image.Id), nil)
	if err != nil {
		panic(err)
	}

	resp, err := wac.conn.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		wac.logger.OK("Image", image.Id, "removed successful")
	}
	wg.Done()
}

func (wac WorkerAsClient) Sweep() {
	fmt.Println("")

	chContainers := make(chan []Container)
	chImages := make(chan []Image)

	go wac.ListContainers(chContainers)
	go wac.ListImages(chImages)

	select {
	case containers := <-chContainers:
		wac.logger.Info("Time:", time.Now().UnixNano())
		wac.logger.Info("Host:", wac.host)
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
	return fmt.Sprintf("v%v", wac.version)
}

func New(host, version string, logger logger.Logger, removeImage, removeContainersPaused bool) (Worker, error) {
	urlParsed, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	conn := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", urlParsed.Path)
			},
		},
	}

	return WorkerAsClient{
		conn,
		host,
		version,
		logger,
		removeImage,
		removeContainersPaused,
	}, nil
}
