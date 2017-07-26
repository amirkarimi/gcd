package main

import  (
	"context"
	"net"
	"net/url"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"sync"
)

type Container struct {
	Id string
	Created int
	SizeRw int
	Status string
	State string
}

func main() {
	host := "unix:///var/run/docker.sock"
	urlParsed, err := url.Parse(host)
	if err != nil {
		fmt.Println(err.Error())
	}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", urlParsed.Path)
			},
		},
	}

	reqQuery, err := http.NewRequest("GET", "http://v1.24/containers/json?all=1", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	respQuery, err := httpc.Do(reqQuery)
	defer respQuery.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	containers := make([]Container, 0)

	bodyQuery, err := ioutil.ReadAll(respQuery.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	json.Unmarshal(bodyQuery, &containers)

	fmt.Println("Host:", urlParsed.Path)

	fmt.Println("Containers total:", len(containers))
	
	var wg sync.WaitGroup
	for _, container := range containers {
		wg.Add(1)
		go func(c Container) {
			if c.State != "running" {

				reqKill, err := http.NewRequest("DELETE", fmt.Sprintf("http://v1.24/containers/%v", c.Id), nil)
				if err != nil {
					fmt.Println(err.Error())
				}

				respKill, err := httpc.Do(reqKill)
				defer respKill.Body.Close()
				if err != nil {
					fmt.Sprintln(err.Error())
				}

				if respKill.StatusCode == http.StatusNoContent {
					fmt.Println("Container", c.Id, "removed succefull!")
				}
			}
			wg.Done()
		}(container)
	}
	wg.Wait()
}
