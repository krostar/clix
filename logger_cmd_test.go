package clix

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/krostar/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoggerFromContext(t *testing.T) {
	t.Run("with logger in context", func(t *testing.T) {
		noop := new(logger.Logger)
		*noop = &logger.Noop{}
		ctx := context.WithValue(context.Background(), ctxKeyLogger, noop)
		log := LoggerFromContext(ctx)
		require.NotNil(t, log)
		assert.IsType(t, &logger.Noop{}, log)
	})
	t.Run("without logger in context", func(t *testing.T) {
		log := LoggerFromContext(context.Background())
		require.Nil(t, log)
	})
}

func Test_WithLogger(t *testing.T) {
	t.Run("default should be enough to log an info", func(t *testing.T) {
		outputRaw, err := logger.CaptureOutput(func() { // capture everything printed to std{out,err}
			err := Command(WithLogger(func(ctx context.Context) (*cobra.Command, context.Context, error) {
				return &cobra.Command{
					RunE: ExecHandler(ctx, func(help func()) (Handler, error) {
						return HandlerFunc(func(ctx context.Context, _, _ []string) error {
							log := LoggerFromContext(ctx)
							log.WithField("hello", "world").Info("displayed")
							return nil
						}), nil
					}),
				}, ctx, nil
			})).Exec(context.Background(), []string{"-f", "json"})
			assert.NoError(t, err)
		})
		require.NoError(t, err)
		var output map[string]interface{} // only one json log is supposed to be wrote
		require.NoError(t, json.Unmarshal([]byte(outputRaw), &output), outputRaw)
		assert.Empty(t, cmp.Diff(map[string]interface{}{ // make sure we get the right log
			"level": logrus.InfoLevel.String(),
			"hello": "world",
			"msg":   "displayed",
		}, output,
			cmpopts.IgnoreMapEntries(func(key string, _ interface{}) bool {
				return key == "time"
			}),
		))
	})

	t.Run("options should be applied", func(t *testing.T) {
		err := Command(WithLogger(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return &cobra.Command{
				SilenceErrors: true,
				SilenceUsage:  true,
				Run:           func(*cobra.Command, []string) {},
			}, ctx, nil
		}, LoggerWithCreateFunc(func(log logger.Config) (logger.Logger, error) {
			return nil, errors.New("boum")
		}))).Exec(context.Background(), []string{})
		require.Error(t, err)
	})

	t.Run("provided command failed to be built", func(t *testing.T) {
		err := Command(WithLogger(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return nil, nil, errors.New("boum")
		})).Exec(context.Background(), []string{})
		require.Error(t, err)
	})
}

func Test_loggerPreRunInit(t *testing.T) {
	t.Run("default logger and option configuration should be enough", func(t *testing.T) {
		var cfg logger.Config
		cfg.SetDefault()

		cmd := cobra.Command{
			PersistentPreRunE: loggerPreRunInit(
				defaultLoggerCommandOptions().createLoggerFunc,
				&cfg,
				new(logger.Logger),
			),
			SilenceErrors: true,
			SilenceUsage:  true,
			Run:           func(*cobra.Command, []string) {},
		}
		assert.NoError(t, cmd.Execute())
	})

	t.Run("logger configuration is invalid", func(t *testing.T) {
		cmd := cobra.Command{
			PersistentPreRunE: loggerPreRunInit(defaultLoggerCommandOptions().createLoggerFunc, &logger.Config{
				Formatter: "boum",
			}, new(logger.Logger)),
			SilenceErrors: true,
			SilenceUsage:  true,
			Run:           func(*cobra.Command, []string) {},
		}
		assert.Error(t, cmd.Execute())
	})

	t.Run("logger initialization failed", func(t *testing.T) {
		var cfg logger.Config
		cfg.SetDefault()

		cmd := cobra.Command{
			PersistentPreRunE: loggerPreRunInit(func(logger.Config) (logger.Logger, error) {
				return nil, errors.New("boum")
			}, &cfg, new(logger.Logger)),
			SilenceErrors: true,
			SilenceUsage:  true,
			Run:           func(*cobra.Command, []string) {},
		}
		assert.Error(t, cmd.Execute())
	})
}
