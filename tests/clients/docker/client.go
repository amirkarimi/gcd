package docker

import (
	"github.com/stone-payments/gcd/clients/docker"
)

type Client struct {
	FakeGetContainers   func() ([]docker.Container, error)
	FakeGetImages       func() ([]docker.Image, error)
	FakeRemoveContainer func(id string) (bool, error)
	FakeRemoveImage     func(id string) (bool, error)
	FakeGetVersion      func() string
	FakeGetHost         func() string
}

func (c Client) GetContainers() ([]docker.Container, error) {
	if c.FakeGetContainers != nil {
		return c.FakeGetContainers()
	}
	return nil, nil
}

func (c Client) GetImages() ([]docker.Image, error) {
	if c.FakeGetImages != nil {
		return c.FakeGetImages()
	}
	return nil, nil
}

func (c Client) RemoveContainer(id string) (bool, error) {
	if c.FakeRemoveContainer != nil {
		return c.FakeRemoveContainer(id)
	}
	return false, nil
}

func (c Client) RemoveImage(id string) (bool, error) {
	if c.FakeRemoveImage != nil {
		return c.FakeRemoveImage(id)
	}
	return false, nil
}

func (c Client) GetVersion() string {
	if c.FakeGetVersion != nil {
		return c.FakeGetVersion()
	}
	return ""
}

func (c Client) GetHost() string {
	if c.FakeGetHost != nil {
		return c.FakeGetHost()
	}
	return ""
}
