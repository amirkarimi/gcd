package main

import  (
	"strconv"
	"os"
	"os/signal"
	"context"
	"net"
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"sync"
	"time"
	"syscall"
)

type (
	Container struct {
		Id string
		Created int
		SizeRw int
		Status string
		State string
	}

	Image struct {
		Id string
	}
)

func removeImage(wg *sync.WaitGroup, httpc http.Client, i Image) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://v1.24/images/%v?force=true", i.Id), nil)
	if err != nil {
		panic(err)
	}

	resp, err := httpc.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("[INFO] Image", i.Id, "removed successful")
	}
	wg.Done()
}

func getImages(cImages chan []Image, httpc http.Client) {
	
	req, err := http.NewRequest("GET", "http://v1.24/images/json?all=1", nil)
	if err != nil {
		panic(err)
	}

	resp, err := httpc.Do(req)
	defer resp.Body.Close()
	if err != nil {
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

	cImages <- images
}

func removeContainer(wg *sync.WaitGroup, httpc http.Client, c Container)  {
	if c.State != "running" {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("http://v1.24/containers/%v", c.Id), nil)
		if err != nil {
			panic(err)
		}

		resp, err := httpc.Do(req)
		if err != nil {
			fmt.Sprintln(err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			fmt.Println("[INFO] Container", c.Id, "removed successful")
		}
	}
	wg.Done()
}

func getContainers(cContainers chan []Container, httpc http.Client) {
	req, err := http.NewRequest("GET", "http://v1.24/containers/json?all=1", nil)
	if err != nil {
		panic(err)
	}

	resp, err := httpc.Do(req)
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

	cContainers <- containers
}

func newClient(host string) (http.Client, error) {
	urlParsed, err := url.Parse(host)
	if err != nil {
		return http.Client{}, err
	}

	return http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", urlParsed.Path)
			},
		},
	}, nil
}

func start(host string, httpc http.Client) {
	cContainers := make(chan []Container)
	cImages := make(chan []Image)

	go getContainers(cContainers, httpc)
	go getImages(cImages, httpc)

	select {
		// case containers, images := <- cContainers, <- cImages:
		case containers := <- cContainers:
			fmt.Println("[INFO] Time:", time.Now().UnixNano())
			fmt.Println("[INFO] Host:", host)
			fmt.Println("[INFO] Containers total:", len(containers))
			select {
				case images := <- cImages:
					fmt.Println("[INFO] Images total:", len(images))

					wgContainerDone := make(chan bool)
					var wgContainer sync.WaitGroup
					for _, container := range containers {
						wgContainer.Add(1)
						go removeContainer(&wgContainer, httpc, container)
					}

					go func(done chan bool) {
						wgContainer.Wait()
						done <- true
					}(wgContainerDone)

					select {
						case <-wgContainerDone:
							fmt.Println("[INFO] Action to containers finished")
						case <-time.After(time.Second * 1):
							fmt.Println("[ERROR] Action to containers timeout occurred")
					}

					wgImageDone := make(chan bool)
					var wgImage sync.WaitGroup
					for _, image := range images {
						wgImage.Add(1)
						go removeImage(&wgImage, httpc, image)
					}
					go func(done chan bool) {
						wgImage.Wait()
						done <- true
					}(wgImageDone)

					select {
						case <-wgImageDone:
							fmt.Println("[INFO] Action to images finished")
						case <-time.After(time.Second * 1):
							fmt.Println("[ERROR] Action to images timeout occurred")
					}

				case <-time.After(time.Second * 1):
					fmt.Println("[ERROR] Images request timeout")
			}

		case <-time.After(time.Second * 1):
			fmt.Println("[ERROR] Containers request timeout")
	}
}

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

	httpc, err := newClient(host)
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
				fmt.Println("")
				start(host, httpc)
		}
	}
}
