package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SetupSignalHandler sets up a signal handler to listen for signals
// and cancel the context when a signal is received.
func SetupSignalHandler(ctx context.Context, cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-signals:
			cancel()
		case <-ctx.Done():
		}
	}()
}
