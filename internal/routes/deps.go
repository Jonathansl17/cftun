package routes

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/service"
)

type Store interface {
	Exists() bool
	Load() (*ingress.Document, error)
	Save(ctx context.Context, doc *ingress.Document) error
	Restore(ctx context.Context) error
}

type Tunnels interface {
	ListTunnels(ctx context.Context) ([]cloudflared.Tunnel, error)
	RouteDNS(ctx context.Context, tunnel, hostname string) error
	ValidateIngress(ctx context.Context, path string) (string, error)
}

type Service interface {
	Run(ctx context.Context, action service.Action) error
}

type DNS interface {
	Configured() (bool, error)
	DeleteCNAME(ctx context.Context, hostname string) (int, error)
}

type Reporter interface {
	Printf(format string, args ...any)
}

type Options struct {
	SkipDNS     bool
	SkipRestart bool
}
