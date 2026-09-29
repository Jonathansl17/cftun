package tunnelcmd

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/cloudflared"
)

type Tunnels interface {
	ListTunnels(ctx context.Context) ([]cloudflared.Tunnel, error)
	CreateTunnel(ctx context.Context, name string) error
	DeleteTunnel(ctx context.Context, name string) error
}

type Initializer interface {
	Init(ctx context.Context, tunnelRef string) (string, error)
}

type Config interface {
	Path() string
	Exists() bool
}

type Deps struct {
	Session     uikit.Session
	Tunnels     Tunnels
	Initializer Initializer
	Config      Config
}

type Group struct {
	deps Deps
}
