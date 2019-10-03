package clix

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_initContextCancelableBySignal(t *testing.T) {
	t.Run("calling cancel func cancels the context", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		assert.NoError(t, ctx.Err())
		cancel()
		<-ctx.Done()
		assert.Error(t, ctx.Err())
	})

	t.Run("sending provided signal cancels the context", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		defer cancel()
		assert.NoError(t, ctx.Err())
		assert.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR1))
		<-ctx.Done()
		assert.Error(t, ctx.Err())
	})

	t.Run("sending unknown signal keeps context intact", func(t *testing.T) {
		ctx, cancel := NewContextCancelableBySignal(syscall.SIGUSR1)
		defer cancel()
		assert.NoError(t, ctx.Err())
		assert.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR2))
		assert.NoError(t, ctx.Err())
	})
}
