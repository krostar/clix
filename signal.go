package clix

import (
	"context"
	"os"
	"os/signal"
)

// NewContextCancelableBySignal creates a new context that cancels
// when provided signals are triggered.
func NewContextCancelableBySignal(signals ...os.Signal) (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	clean := func() {
		signal.Ignore(signals...)
		close(signalChan)
	}

	// catch some stop signals, and cancel the context if caught
	signal.Notify(signalChan, signals...)
	go func() {
		<-signalChan // block until a signal is received
		cancel()
	}()

	return ctx, clean
}
