package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Jonathansl17/cftun/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), shutdownSignals...)
	code := cli.Run(ctx, cli.Streams{In: os.Stdin, Out: os.Stdout, Err: os.Stderr}, version)
	stop()
	os.Exit(code)
}
