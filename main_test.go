package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
)

func TestRemoveHealthyContainersExited(t *testing.T) {
	// Run a sample docker
	runShell(t, "docker run --name hello-test hello-world")

	dc, err := docker.NewClient(options.dockerHost)
	if err != nil {
		t.Error(err)
	}

	Run(
		dc,
		Options{
			dockerHost:                    "unix:///var/run/docker.sock",
			removeImages:                  false,
			removeHealthyContainersExited: true,
		},
	)

	// Check if the sample docker removed
	assertContainerNotExists(t, dc, "hello-test")
}

func TestRemoveImages(t *testing.T) {
	// Run a sample docker
	runShell(t, "docker pull hello-world")

	dc, err := docker.NewClient(options.dockerHost)
	if err != nil {
		t.Error(err)
	}

	Run(
		dc,
		Options{
			dockerHost:                    "unix:///var/run/docker.sock",
			removeImages:                  true,
			removeHealthyContainersExited: false,
		},
	)

	// Check if the sample docker removed
	assertImageNotExists(t, dc, "hello-world")
}

func assertContainerNotExists(t *testing.T, dc *docker.Client, name string) {
	containers, err := dc.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	if err != nil {
		t.Error(err)
	}
	for _, c := range containers {
		if c.Names[0] == "hello-test" {
			t.Errorf("Container %s exists which shouldn't.", name)
			return
		}
	}
}

func assertImageNotExists(t *testing.T, dc *docker.Client, tag string) {
	images, err := dc.ListImages(docker.ListImagesOptions{
		All: true,
	})
	if err != nil {
		t.Error(err)
	}
	for _, i := range images {
		for _, tag := range i.RepoTags {
			if tag == "hello-world:latest" {
				t.Errorf("Image %s exists which shouldn't.", tag)
			}
		}
	}
}

func runShell(t *testing.T, cmdFull string) string {
	fullParams := strings.Split(cmdFull, " ")
	cmd := fullParams[0]
	params := fullParams[1:]
	stdout, err := exec.Command(cmd, params...).Output()
	if err != nil {
		t.Error(err)
	}
	return fmt.Sprintf("%s", stdout)
}
