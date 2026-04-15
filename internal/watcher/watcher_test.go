package watcher_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logdrift/internal/watcher"
)

func TestNew_DefaultPollInterval(t *testing.T) {
	w := watcher.New("/dev/null", 0)
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}

func TestRun_FileNotFound(t *testing.T) {
	w := watcher.New("/nonexistent/path/logdrift.log", time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRun_EmitsNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logdrift-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	w := watcher.New(f.Name(), 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Run(ctx)
	}()

	// Give the watcher time to seek to EOF before we write.
	time.Sleep(30 * time.Millisecond)

	want := "INFO hello world"
	if _, err := f.WriteString(want + "\n"); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case got := <-w.Lines:
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for line")
	}

	cancel()
	<-errCh // drain
}

func TestRun_ClosesChannelOnCancel(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logdrift-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	w := watcher.New(f.Name(), 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	go w.Run(ctx) //nolint:errcheck
	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case _, ok := <-w.Lines:
		if ok {
			t.Error("expected Lines channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel close")
	}
}
