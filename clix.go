// Package clix ...
package clix

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// CLI stores the command line interfaces command and subcommands.
type CLI struct {
	command     CommandBuilderFunc
	subcommands []CommandBuilderFunc
}

// CommandBuilderFunc defines a cobra command builder func.
type CommandBuilderFunc func(context.Context) (*cobra.Command, context.Context, error)

// Command creates a new root cli instance.
func Command(command CommandBuilderFunc) *CLI { return &CLI{command: command} }

// SubCommand adds a sub command to the main command.
func (cli *CLI) SubCommand(cmd CommandBuilderFunc) *CLI {
	cli.subcommands = append(cli.subcommands, cmd)
	return cli
}

// Build return a concatenated command builder that adds all subcommands to the root command.
func (cli *CLI) Build() CommandBuilderFunc {
	return func(ctx context.Context) (*cobra.Command, context.Context, error) {
		command, ctx, err := cli.command(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to build command: %w", err)
		}
		for _, subcmd := range cli.subcommands {
			sub, _, err := subcmd(ctx)
			if err != nil {
				return nil, nil, fmt.Errorf("unable to build subcommand: %w", err)
			}
			command.AddCommand(sub)
		}
		return command, ctx, nil
	}
}

// Exec executes the command given by args.
func (cli *CLI) Exec(ctx context.Context, args []string) error {
	cmd, _, err := cli.Build()(ctx)
	if err != nil {
		return fmt.Errorf("unable to build command: %w", err)
	}
	cmd.SetArgs(args)
	return cmd.Execute()
}

type (
	// GetHandlerFunc returns a handler, providing it's help function.
	GetHandlerFunc func(help func()) (Handler, error)
	// Handler abstracts a cli command handler.
	Handler interface {
		Handle(ctx context.Context, args []string, dashedArgs []string) error
	}
	// HandlerFunc implements Handler.
	HandlerFunc func(ctx context.Context, args []string, dashedArgs []string) error
)

// Handle implement Handler.
func (h HandlerFunc) Handle(ctx context.Context, args []string, dashedArgs []string) error {
	return h(ctx, args, dashedArgs)
}

// ExecHandler execs the provided handler function.
func ExecHandler(ctx context.Context, getHandler GetHandlerFunc) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, args []string) error {
		// before reaching this point, we want the usage but after that, if we
		// have an error, we do want to handle when to display it, and when not
		c.SilenceUsage = true

		cmd, err := getHandler(func() {
			c.Help() // nolint: errcheck, gosec
		})
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return cmd.Handle(ctx, nil, nil)
		}

		argsSeparatedAt := c.ArgsLenAtDash()
		switch {
		case argsSeparatedAt == 0 && len(args) > 0:
			return cmd.Handle(ctx, nil, args)
		case argsSeparatedAt > 0 && len(args[argsSeparatedAt:]) > 0:
			return cmd.Handle(ctx, args[:argsSeparatedAt], args[argsSeparatedAt:])
		default:
			return cmd.Handle(ctx, args, nil)
		}
	}
}
