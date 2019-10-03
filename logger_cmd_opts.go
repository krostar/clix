package clix

import (
	"github.com/krostar/logger"
	"github.com/krostar/logger/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type loggerCommandOptions struct {
	appName            string
	appVersion         string
	createLoggerFunc   func(cfg logger.Config) (logger.Logger, error)
	setPersistentFlags func(flags *pflag.FlagSet, cfg *logger.Config)
}

func (c *loggerCommandOptions) applyToCommand(cmd *cobra.Command) {
	if c.appName != "" {
		cmd.Use = c.appName
	}
	if c.appVersion != "" {
		cmd.Version = c.appVersion
	}
}

func defaultLoggerCommandOptions() *loggerCommandOptions {
	return &loggerCommandOptions{
		appName:    "cli-app",
		appVersion: "dev",
		createLoggerFunc: func(cfg logger.Config) (logger.Logger, error) {
			return logrus.New(logrus.WithConfig(cfg))
		},
		setPersistentFlags: func(flags *pflag.FlagSet, cfg *logger.Config) {
			flags.StringVarP(&cfg.Verbosity,
				"log-verbosity", "v", cfg.Verbosity,
				"verbosity of logs printed to the standard output",
			)
			flags.StringVarP(&cfg.Formatter,
				"log-format", "f", cfg.Formatter,
				"format to print logs to standard output with",
			)
		},
	}
}

// LoggerCommandOption defines the signature of an option applier.
type LoggerCommandOption func(o *loggerCommandOptions)

// LoggerWithAppName sets the root command app name.
func LoggerWithAppName(appName string) LoggerCommandOption {
	return func(o *loggerCommandOptions) { o.appName = appName }
}

// LoggerWithVersion sets the root command app version.
func LoggerWithVersion(version string) LoggerCommandOption {
	return func(o *loggerCommandOptions) { o.appVersion = version }
}

// LoggerWithPersistentFlagsFunc overrides the defaults persistent flag set.
func LoggerWithPersistentFlagsFunc(fct func(flags *pflag.FlagSet, cfg *logger.Config)) LoggerCommandOption {
	return func(o *loggerCommandOptions) { o.setPersistentFlags = fct }
}

// LoggerWithCreateFunc sets the logger creation function.
func LoggerWithCreateFunc(fct func(log logger.Config) (logger.Logger, error)) LoggerCommandOption {
	return func(o *loggerCommandOptions) { o.createLoggerFunc = fct }
}
