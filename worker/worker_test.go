package worker

import (
	"errors"
	"reflect"
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
			t.Error("Timeout test GetContainers")
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
			t.Error("Timeout test GetImages")
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
		}
	}
}
