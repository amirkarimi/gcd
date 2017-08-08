package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
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

	Client interface {
		GetContainers() ([]Container, error)
		GetImages() ([]Image, error)
		RemoveContainer(id string) (bool, error)
		RemoveImage(id string) error
		GetVersion() string
		GetHost() string
	}

	DockerAsClient struct {
		host    string
		conn    http.Client
		version string
	}
)

func (dac DockerAsClient) GetContainers() ([]Container, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://v%v/containers/json?all=1", dac.version), nil)
	if err != nil {
		return nil, err
	}

	resp, err := dac.conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	containers := make([]Container, 0)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &containers); err != nil {
		return nil, err
	}

	return containers, nil
}

func (dac DockerAsClient) GetImages() ([]Image, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://v%v/images/json?all=1", dac.version), nil)
	if err != nil {
		return nil, err
	}

	resp, err := dac.conn.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	images := make([]Image, 0)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &images); err != nil {
		return nil, err
	}

	return images, nil
}

func (dac DockerAsClient) RemoveContainer(id string) (bool, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://v%v/containers/%v", dac.version, id), nil)
	if err != nil {
		return false, err
	}

	resp, err := dac.conn.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	}

	return false, nil
}

func (dac DockerAsClient) RemoveImage(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://v%v/images/%v?force=true", dac.version, id), nil)
	if err != nil {
		return err
	}

	resp, err := dac.conn.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (dac DockerAsClient) GetVersion() string {
	return dac.version
}

func (dac DockerAsClient) GetHost() string {
	return dac.host
}

func NewClient(host, version string) (Client, error) {
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

	return DockerAsClient{
		host,
		conn,
		version,
	}, nil
}
