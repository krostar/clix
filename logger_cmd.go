package clix

import (
	"context"
	"fmt"

	"github.com/krostar/logger"
	"github.com/spf13/cobra"
)

type ctxKey string

const ctxKeyLogger ctxKey = "logger"

// LoggerFromContext returns the logger from the context, if present
func LoggerFromContext(ctx context.Context) logger.Logger {
	if logPtr, hasLogger := ctx.Value(ctxKeyLogger).(*logger.Logger); hasLogger && logPtr != nil {
		return *logPtr
	}
	return nil
}

// WithLogger adds to an existing command log flags, and config requirements.
func WithLogger(cbf CommandBuilderFunc, opts ...LoggerCommandOption) CommandBuilderFunc {
	return func(ctx context.Context) (*cobra.Command, context.Context, error) {
		log := new(logger.Logger)
		ctx = context.WithValue(ctx, ctxKeyLogger, log)

		cmd, ctx, err := cbf(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to build root command: %w", err)
		}

		o := defaultLoggerCommandOptions()
		for _, opt := range opts {
			opt(o)
		}
		o.applyToCommand(cmd)

		var cfg logger.Config
		cfg.SetDefault()

		o.setPersistentFlags(cmd.PersistentFlags(), &cfg)
		cmd.PersistentPreRunE = loggerPreRunInit(o.createLoggerFunc, &cfg, log)

		return cmd, ctx, nil
	}
}

func loggerPreRunInit(
	createLogger func(cfg logger.Config) (logger.Logger, error),
	cfg *logger.Config,
	logPtr *logger.Logger,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("logger config is invalid: %w", err)
		}

		var err error
		*logPtr, err = createLogger(*cfg)
		if err != nil {
			return fmt.Errorf("unable to create logger: %w", err)
		}

		return nil
	}
}
