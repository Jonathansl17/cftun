package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jonathansl17/cftun/internal/cli"
	"github.com/Jonathansl17/cftun/internal/msg"
)

const exitFailure = 1

var version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	err := cli.NewRoot(version, os.Stdin, os.Stdout).ExecuteContext(ctx)
	stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, msg.ErrorFormat, err)
		os.Exit(exitFailure)
	}
}
