package clix

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLI_Exec(t *testing.T) {
	t.Run("no command should display help without error", func(t *testing.T) {
		cli := Command(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return &cobra.Command{Run: func(*cobra.Command, []string) {}}, ctx, nil
		})
		assert.NoError(t, cli.Exec(context.Background(), []string{""}))
	})

	t.Run("sub command exists and should be called", func(t *testing.T) {
		cli := Command(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return &cobra.Command{Run: func(*cobra.Command, []string) {}}, ctx, nil
		})
		called := false
		cli.SubCommand(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return &cobra.Command{
				Use: "sub",
				Run: func(*cobra.Command, []string) { called = true },
			}, ctx, nil
		})
		assert.NoError(t, cli.Exec(context.Background(), []string{"sub"}))
		assert.True(t, called)
	})

	t.Run("build root command failed", func(t *testing.T) {
		cli := Command(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return &cobra.Command{Run: func(*cobra.Command, []string) {}}, ctx, errors.New("boum")
		})
		cli = cli.SubCommand(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return &cobra.Command{
				Use: "sub",
				Run: func(*cobra.Command, []string) {},
			}, ctx, nil
		})
		assert.Error(t, cli.Exec(context.Background(), []string{"sub"}))
	})

	t.Run("build subcommand failed", func(t *testing.T) {
		cli := Command(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return &cobra.Command{Run: func(*cobra.Command, []string) {}}, ctx, nil
		})
		cli = cli.SubCommand(func(ctx context.Context) (*cobra.Command, context.Context, error) {
			return &cobra.Command{
				Use: "sub",
				Run: func(*cobra.Command, []string) {},
			}, ctx, errors.New("boum")
		})
		assert.Error(t, cli.Exec(context.Background(), []string{"sub"}))
	})
}

func Test_ExecHandler_handler_called_without_error(t *testing.T) {
	cmd := &cobra.Command{Use: "sub", SilenceErrors: true}
	cmd.SetOutput(ioutil.Discard)

	called := false
	cmd.RunE = ExecHandler(context.Background(), func(help func()) (Handler, error) {
		return HandlerFunc(func(ctx context.Context, args, dashed []string) error {
			assert.Equal(t, []string{"a", "b"}, args)
			assert.Equal(t, []string{"c", "d"}, dashed)
			assert.NotNil(t, ctx)
			assert.NotPanics(t, help)
			assert.True(t, cmd.SilenceUsage)
			called = true
			return nil
		}), nil
	})
	cmd.SetArgs([]string{"a", "b", "--", "c", "d"})
	require.NoError(t, cmd.Execute())
	assert.True(t, called)
}

func Test_ExecHandler_called_handler_returned_in_error(t *testing.T) {
	cmd := &cobra.Command{Use: "sub", SilenceErrors: true}
	cmd.RunE = ExecHandler(context.Background(), func(help func()) (Handler, error) {
		return HandlerFunc(func(ctx context.Context, args, dashed []string) error {
			return errors.New("boum")
		}), nil
	})
	require.Error(t, cmd.Execute())
}

func Test_ExecHandler_handler_getter_failed(t *testing.T) {
	cmd := &cobra.Command{Use: "sub", SilenceErrors: true}
	cmd.RunE = ExecHandler(context.Background(), func(help func()) (Handler, error) {
		return nil, errors.New("boum")
	})
	require.Error(t, cmd.Execute())
}

func Test_ExecHandler_arg_and_dash_args_are_nil_when_handler_is_called_without_args(t *testing.T) {
	cmd := &cobra.Command{Use: "sub", SilenceErrors: true}
	cmd.RunE = ExecHandler(context.Background(), func(help func()) (Handler, error) {
		return HandlerFunc(func(ctx context.Context, args, dashed []string) error {
			assert.Nil(t, args)
			assert.Nil(t, dashed)
			return nil
		}), nil
	})
	cmd.SetArgs([]string{})
	require.NoError(t, cmd.Execute())
}

func Test_ExecHandler_arg_and_dash_args_are_nil_when_handler_is_called_only_with_double_dash(t *testing.T) {
	cmd := &cobra.Command{Use: "sub", SilenceErrors: true}
	cmd.RunE = ExecHandler(context.Background(), func(help func()) (Handler, error) {
		return HandlerFunc(func(ctx context.Context, args, dashed []string) error {
			assert.Nil(t, args)
			assert.Nil(t, dashed)
			return nil
		}), nil
	})
	cmd.SetArgs([]string{"--"})
	require.NoError(t, cmd.Execute())
}

func Test_ExecHandler_arg_are_nil_when_handler_is_called_only_with_dash_args(t *testing.T) {
	cmd := &cobra.Command{Use: "sub", SilenceErrors: true}
	cmd.RunE = ExecHandler(context.Background(), func(help func()) (Handler, error) {
		return HandlerFunc(func(ctx context.Context, args, dashed []string) error {
			assert.Nil(t, args)
			assert.Equal(t, []string{"c", "d"}, dashed)
			return nil
		}), nil
	})
	cmd.SetArgs([]string{"--", "c", "d"})
	require.NoError(t, cmd.Execute())
}

func Test_ExecHandler_dashed_arg_are_nil_when_handler_is_called_only_with_args(t *testing.T) {
	cmd := &cobra.Command{Use: "sub", SilenceErrors: true}
	cmd.RunE = ExecHandler(context.Background(), func(help func()) (Handler, error) {
		return HandlerFunc(func(ctx context.Context, args, dashed []string) error {
			assert.Equal(t, []string{"a", "b"}, args)
			assert.Nil(t, dashed)
			return nil
		}), nil
	})
	cmd.SetArgs([]string{"a", "b"})
	require.NoError(t, cmd.Execute())
}

func Test_ExecHandler_dashed_arg_are_nil_when_handler_is_called_with_nothing_after_doule_dash(t *testing.T) {
	cmd := &cobra.Command{Use: "sub", SilenceErrors: true}
	cmd.RunE = ExecHandler(context.Background(), func(help func()) (Handler, error) {
		return HandlerFunc(func(ctx context.Context, args, dashed []string) error {
			assert.Equal(t, []string{"a", "b"}, args)
			assert.Nil(t, dashed)
			return nil
		}), nil
	})
	cmd.SetArgs([]string{"a", "b", "--"})
	require.NoError(t, cmd.Execute())
}
