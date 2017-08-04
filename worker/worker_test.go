package worker

import (
	"os"
	"strconv"
	"testing"
)

var (
	sock    = "/var/run/docker.sock"
	version = "1.24"
)

func TestListContainers(t *testing.T) {
	worker, err := New(sock, version)
	if err != nil {
		panic(err)
	}

	ch := make(chan []Container)

	go worker.ListContainers(ch)

	containers := <-ch

	result := len(containers)
	expected, err := strconv.Atoi(os.Getenv("GCD_TEST_AMOUNT_CONTAINERS"))
	if err != nil {
		panic(err)
	}

	if result != expected {
		t.Errorf("worker.ListContainers, result: %v, expected: %v", result, expected)
	}
}

func BenchmarkListContainers(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	worker, err := New(sock, version)
	if err != nil {
		panic(err)
	}

	ch := make(chan []Container)

	for i := 0; i < b.N; i++ {
		go worker.ListContainers(ch)
	}
}
