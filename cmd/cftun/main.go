// Command cftun manages a Cloudflare Tunnel and its public hostnames.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jonathansl17/cloudfare/internal/cli"
	"github.com/Jonathansl17/cloudfare/internal/msg"
)

const exitFailure = 1

// version is set at build time with -ldflags "-X main.version=vX.Y.Z".
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
