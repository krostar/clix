package clix_test

import (
	"context"
	"fmt"
	"log"

	"github.com/krostar/clix"
	"github.com/krostar/logger"
	"github.com/spf13/cobra"
)

func Example() {
	// cli is buildable and modulable
	// here we use the default root command that does nothing
	cli := clix.Command(commandA)

	// or even more complex command tree
	cli = cli.SubCommand(clix.
		Command(commandB).
		SubCommand(commandBB).
		Build(),
	)

	// let's add even more commands
	cli = cli.SubCommand(commandC)

	/*
		Which produces the following command tree:

		commandA
		├── commandB
		│   └── commandBB
		└── commandC
	*/

	// once the cli is defined, we only have to execute it with chosen context and arguments
	// context can be set to be canceled by CTRL+C, timeout, ...
	// depending on the task executed (cron, daemon, ...)
	// here we use an empty one
	ctx := context.Background()
	// args should probably be set to os.Args
	// but for simplicity sake here we hardcode them
	args := []string{"commandB", "commandBB"}

	if err := cli.Exec(ctx, args); err != nil {
		// something didn't worked as expected
		log.Println(fmt.Errorf("cli execution failed: %w", err))
	}

	// Output:
	// bb handled
}

// commandBB builds the BB cobra command, define flags, ...
func commandBB(ctx context.Context) (*cobra.Command, context.Context, error) {
	return &cobra.Command{
		Use: "commandBB",
		RunE: clix.ExecHandler(ctx, func(showHelp func()) (clix.Handler, error) {
			// inject any dependencies the bb command handler needs
			return &handleBB{
				help: showHelp,
				log:  clix.LoggerFromContext(ctx),
			}, nil
		}),
	}, ctx, nil
}

// handleBB implements clix.Handler, is a clean, cobra-unrelated, reusable
// function that handle the bb command
type handleBB struct {
	help func()
	log  logger.Logger
}

// abstract cli handler that call whatener business logic it needs
// following whatever architecture it wants, 100% testable as every dependencies are injected
func (h handleBB) Handle(ctx context.Context, args, dashed []string) error {
	fmt.Println("bb handled")
	return nil
}

func commandA(ctx context.Context) (*cobra.Command, context.Context, error) {
	return &cobra.Command{
		Use: "commandA",
		Run: func(cmd *cobra.Command, args []string) {},
	}, ctx, nil
}

func commandB(ctx context.Context) (*cobra.Command, context.Context, error) {
	return &cobra.Command{
		Use: "commandB",
		Run: func(cmd *cobra.Command, args []string) {},
	}, ctx, nil
}

func commandC(ctx context.Context) (*cobra.Command, context.Context, error) {
	return &cobra.Command{
		Use: "commandC",
		Run: func(cmd *cobra.Command, args []string) {},
	}, ctx, nil
}
