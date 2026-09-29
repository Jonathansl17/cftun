package cli

import (
	"context"
	"fmt"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func Run(ctx context.Context, streams Streams, version string) int {
	root, err := NewRoot(version, streams)
	if err == nil {
		err = root.ExecuteContext(ctx)
	}
	if err != nil {
		fmt.Fprintf(streams.Err, msg.ErrorFormat, uikit.Describe(err))
		return exitFailure
	}
	return exitOK
}
