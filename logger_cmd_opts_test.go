package clix

import (
	"testing"

	"github.com/krostar/logger"
	"github.com/krostar/logger/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_defaultLoggerCommandOptions(t *testing.T) {
	o := defaultLoggerCommandOptions()
	assert.Equal(t, "cli-app", o.appName)
	assert.Equal(t, "dev", o.appVersion)

	t.Run("logrus logger by default", func(t *testing.T) {
		var cfg logger.Config
		cfg.SetDefault()
		log, err := o.createLoggerFunc(cfg)
		require.NoError(t, err)
		assert.IsType(t, &logrus.Logrus{}, log)
	})

	t.Run("long flags should set the config", func(t *testing.T) {
		var cfg logger.Config
		flags := pflag.NewFlagSet("", pflag.ContinueOnError)
		o.setPersistentFlags(flags, &cfg)
		err := flags.Parse([]string{"--log-verbosity", "error", "--log-format", "json"})
		require.NoError(t, err)
		assert.Equal(t, "error", cfg.Verbosity)
		assert.Equal(t, "json", cfg.Formatter)
	})

	t.Run("short flags should set the config", func(t *testing.T) {
		var cfg logger.Config
		flags := pflag.NewFlagSet("", pflag.ContinueOnError)
		o.setPersistentFlags(flags, &cfg)
		err := flags.Parse([]string{"-v", "error", "-f", "json"})
		require.NoError(t, err)
		assert.Equal(t, "error", cfg.Verbosity)
		assert.Equal(t, "json", cfg.Formatter)
	})
}

func TestLoggerCommandOptions_applyToCommand(t *testing.T) {
	var (
		cmd cobra.Command
		o   loggerCommandOptions
	)

	o.applyToCommand(&cmd)
	assert.Equal(t, cobra.Command{}, cmd)

	o.appName = "app"
	o.appVersion = "version"
	o.applyToCommand(&cmd)
	assert.Equal(t, cobra.Command{
		Use:     "app",
		Version: "version",
	}, cmd)
}

func Test_LoggerWithAppName(t *testing.T) {
	var o loggerCommandOptions
	LoggerWithAppName("go-app")(&o)
	assert.Equal(t, "go-app", o.appName)
}

func Test_LoggerWithVersion(t *testing.T) {
	var o loggerCommandOptions
	LoggerWithVersion("go-version")(&o)
	assert.Equal(t, "go-version", o.appVersion)
}

func Test_LoggerWithLoggerCreateFunc(t *testing.T) {
	var o loggerCommandOptions
	LoggerWithCreateFunc(func(logger.Config) (logger.Logger, error) { return nil, nil })(&o)
	assert.NotNil(t, o.createLoggerFunc)
}

func Test_LoggerWithPersistentFlagsFunc(t *testing.T) {
	var o loggerCommandOptions
	LoggerWithPersistentFlagsFunc(func(flags *pflag.FlagSet, cfg *logger.Config) {})(&o)
	assert.NotNil(t, o.setPersistentFlags)
}
