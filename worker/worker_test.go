package worker

import (
	"errors"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stone-payments/gcd/clients/docker"
	"github.com/stone-payments/gcd/logger"
	dockertest "github.com/stone-payments/gcd/tests/clients/docker"
	loggertest "github.com/stone-payments/gcd/tests/logger"
)

func TestListContainers(t *testing.T) {

	cases := []struct {
		Expected          []docker.Container
		FakeGetContainers func() ([]docker.Container, error)
	}{
		{
			Expected: make([]docker.Container, 0),
			FakeGetContainers: func() ([]docker.Container, error) {
				return make([]docker.Container, 0), nil
			},
		},
		{
			Expected: make([]docker.Container, 5),
			FakeGetContainers: func() ([]docker.Container, error) {
				return make([]docker.Container, 5), nil
			},
		},
		{
			Expected: nil,
			FakeGetContainers: func() ([]docker.Container, error) {
				return nil, nil
			},
		},
		{
			Expected: nil,
			FakeGetContainers: func() ([]docker.Container, error) {
				return nil, errors.New("Fake error")
			},
		},
	}

	for _, test := range cases {

		dockerClient := dockertest.Client{
			FakeGetContainers: test.FakeGetContainers,
		}

		w := New(dockerClient, loggertest.Logger{}, false, false)

		ch := make(chan []docker.Container)
		go w.ListContainers(ch)

		select {
		case containers := <-ch:
			if !reflect.DeepEqual(test.Expected, containers) {
				t.Errorf("Not equals containers result: expected: %v, result: %v", test.Expected, containers)
			}
		case <-time.After(time.Second * 4):
			t.Error("Timeout test ListContainers")
		}

	}

}

func BenchmarkListContainers(b *testing.B) {

	b.ResetTimer()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		dockerClient := dockertest.Client{
			FakeGetContainers: func() ([]docker.Container, error) {
				return make([]docker.Container, 0), nil
			},
		}

		w := New(dockerClient, logger.New(), false, false)

		ch := make(chan []docker.Container)
		go w.ListContainers(ch)

		select {
		case <-ch:
		case <-time.After(time.Second * 4):
			log.Println("Goroutine timeout - ListContainers")
		}
	}
}

func TestListImages(t *testing.T) {

	cases := []struct {
		Expected      []docker.Image
		FakeGetImages func() ([]docker.Image, error)
	}{
		{
			Expected: make([]docker.Image, 0),
			FakeGetImages: func() ([]docker.Image, error) {
				return make([]docker.Image, 0), nil
			},
		},
		{
			Expected: make([]docker.Image, 5),
			FakeGetImages: func() ([]docker.Image, error) {
				return make([]docker.Image, 5), nil
			},
		},
		{
			Expected: nil,
			FakeGetImages: func() ([]docker.Image, error) {
				return nil, nil
			},
		},
		{
			Expected: nil,
			FakeGetImages: func() ([]docker.Image, error) {
				return nil, errors.New("Fake error")
			},
		},
	}

	for _, test := range cases {

		dockerClient := dockertest.Client{
			FakeGetImages: test.FakeGetImages,
		}

		w := New(dockerClient, loggertest.Logger{}, false, false)

		ch := make(chan []docker.Image)
		go w.ListImages(ch)

		select {
		case images := <-ch:
			if !reflect.DeepEqual(test.Expected, images) {
				t.Errorf("Not equals images result: expected: %v, result: %v", test.Expected, images)
			}
		case <-time.After(time.Second * 4):
			t.Error("Timeout test ListImages")
		}

	}

}

func BenchmarkListImages(b *testing.B) {

	b.ResetTimer()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		dockerClient := dockertest.Client{
			FakeGetImages: func() ([]docker.Image, error) {
				return make([]docker.Image, 0), nil
			},
		}

		w := New(dockerClient, logger.New(), false, false)

		ch := make(chan []docker.Image)
		go w.ListImages(ch)

		select {
		case <-ch:
		case <-time.After(time.Second * 4):
			log.Println("Goroutine timeout - ListImages")
		}
	}
}

func TestRemoveImage(t *testing.T) {

	cases := []struct {
		Image           docker.Image
		FakeRemoveImage func(id string) error
	}{
		{
			Image: docker.Image{Id: "fake-image-id0"},
			FakeRemoveImage: func(id string) error {
				if id == "fake-image-id0" {
					return nil
				}
				return errors.New("Fake error")
			},
		},
		{
			Image: docker.Image{Id: "fake-image-id1"},
			FakeRemoveImage: func(id string) error {
				if id == "fake-image-id0" {
					return nil
				}
				return errors.New("Fake error")
			},
		},
	}

	for _, test := range cases {

		dockerClient := dockertest.Client{
			FakeRemoveImage: test.FakeRemoveImage,
		}

		w := New(dockerClient, loggertest.Logger{}, false, false)

		var wg sync.WaitGroup

		wg.Add(1)
		w.RemoveImage(&wg, test.Image)

		wg.Wait()
	}

}

func TestRemoveContainer(t *testing.T) {

	FakeRemoveContainer := func(id string) (bool, error) {
		if id == "fake-container-id0" {
			return true, nil
		} else {
			if id == "fake-container-id1" {
				return false, nil
			}
			return false, errors.New("Fake error")
		}
	}

	cases := []struct {
		RemoveContainersExited bool
		Container              docker.Container
		FakeRemoveContainer    func(id string) (bool, error)
	}{
		{
			RemoveContainersExited: false,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "running",
				Status: "Up",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: true,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "running",
				Status: "Up",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: false,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "exited",
				Status: "Exited (0)",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: true,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "exited",
				Status: "Exited (0)",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: false,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "exited",
				Status: "Exited (2)",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: false,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "exited",
				Status: "Exited (2)",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: false,
			Container: docker.Container{
				Id:     "fake-container-id0",
				State:  "exited",
				Status: "Exited (2)",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: true,
			Container: docker.Container{
				Id:     "fake-container-id1",
				State:  "exited",
				Status: "Exited (0)",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
	}

	for _, test := range cases {

		dockerClient := dockertest.Client{
			FakeRemoveContainer: test.FakeRemoveContainer,
		}

		w := New(dockerClient, loggertest.Logger{}, false, test.RemoveContainersExited)

		var wg sync.WaitGroup

		wg.Add(1)
		w.RemoveContainer(&wg, test.Container)

		wg.Wait()
	}
}

func TestGetVersion(t *testing.T) {

	version := "fake-version"

	dockerClient, _ := docker.NewClient("fake-host", version)

	w := New(dockerClient, nil, false, false)

	if w.GetVersion() != version {
		t.Errorf("Unexpected value as return, result: %v, expected: %v", dockerClient.GetVersion(), version)
	}

}

func TestSweep(t *testing.T) {

	FakeRemoveContainer := func(id string) (bool, error) {
		if id == "fake-container-id0" {
			return true, nil
		} else {
			if id == "fake-container-id1" {
				return false, nil
			}
			return false, errors.New("Fake error")
		}
	}

	cases := []struct {
		RemoveContainersExited bool
		RemoveImages           bool
		Container              docker.Container
		FakeRemoveContainer    func(id string) (bool, error)
		FakeGetContainers      func() ([]docker.Container, error)
		FakeGetImages          func() ([]docker.Image, error)
	}{
		{
			RemoveContainersExited: false,
			RemoveImages:           false,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "running",
				Status: "Up",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: false,
			RemoveImages:           true,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "running",
				Status: "Up",
			},
			FakeRemoveContainer: FakeRemoveContainer,
		},
		{
			RemoveContainersExited: false,
			RemoveImages:           true,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "running",
				Status: "Up",
			},
			FakeRemoveContainer: FakeRemoveContainer,
			FakeGetContainers: func() ([]docker.Container, error) {
				return make([]docker.Container, 2), nil
			},
		},
		{
			RemoveContainersExited: false,
			RemoveImages:           true,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "running",
				Status: "Up",
			},
			FakeRemoveContainer: FakeRemoveContainer,
			FakeGetImages: func() ([]docker.Image, error) {
				return make([]docker.Image, 2), nil
			},
		},
		{
			RemoveContainersExited: false,
			RemoveImages:           true,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "running",
				Status: "Up",
			},
			FakeRemoveContainer: FakeRemoveContainer,
			FakeGetContainers: func() ([]docker.Container, error) {
				time.Sleep(time.Second * 5)
				return make([]docker.Container, 2), nil
			},
		},
		{
			RemoveContainersExited: false,
			RemoveImages:           true,
			Container: docker.Container{
				Id:     "fake-container-id-1",
				State:  "running",
				Status: "Up",
			},
			FakeRemoveContainer: FakeRemoveContainer,
			FakeGetImages: func() ([]docker.Image, error) {
				time.Sleep(time.Second * 5)
				return make([]docker.Image, 2), nil
			},
		},
	}

	for _, test := range cases {

		dockerClient := dockertest.Client{
			FakeRemoveContainer: test.FakeRemoveContainer,
			FakeGetContainers:   test.FakeGetContainers,
			FakeGetImages:       test.FakeGetImages,
		}

		w := New(dockerClient, loggertest.Logger{}, test.RemoveImages, test.RemoveContainersExited)

		w.Sweep()
	}

}
