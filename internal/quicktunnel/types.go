package quicktunnel

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/hostfs"
)

type Tunnels interface {
	QuickTunnel(ctx context.Context, url, emptyConfig string) error
}

type TempFiles interface {
	EmptyTemp(pattern string) (hostfs.Temp, error)
}

type Runner struct {
	Tunnels Tunnels
	Files   TempFiles
}
