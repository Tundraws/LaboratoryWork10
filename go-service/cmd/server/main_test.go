package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestRunGracefulShutdown(t *testing.T) {
	stop := make(chan os.Signal, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		stop <- syscall.SIGTERM
	}()

	if err := run("test-secret", "18080", stop); err != nil {
		t.Fatalf("run() error = %v", err)
	}
}

func TestRunGracefulShutdownWithDefaults(t *testing.T) {
	stop := make(chan os.Signal, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		stop <- syscall.SIGINT
	}()

	if err := run("", "", stop); err != nil {
		t.Fatalf("run() with defaults error = %v", err)
	}
}

func TestRunReturnsServerError(t *testing.T) {
	stop := make(chan os.Signal, 1)

	errCh := make(chan error, 1)
	go func() {
		errCh <- run("test-secret", "-1", stop)
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("run() error = nil, want non-nil")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("run() did not return in time")
	}
}

func TestMakeSignalChannel(t *testing.T) {
	stop := makeSignalChannel()
	if stop == nil {
		t.Fatal("makeSignalChannel() returned nil channel")
	}
}
