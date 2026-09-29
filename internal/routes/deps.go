// Package routes adds, edits and removes public hostnames end to end: config
// file, validation, DNS record and service restart.
package routes

import (
	"context"

	"github.com/Jonathansl17/cloudfare/internal/cloudflared"
	"github.com/Jonathansl17/cloudfare/internal/ingress"
	"github.com/Jonathansl17/cloudfare/internal/service"
)

// Store loads and saves the config.
type Store interface {
	Exists() bool
	Load() (*ingress.Document, error)
	Save(ctx context.Context, doc *ingress.Document) error
	Restore(ctx context.Context) error
}

// Tunnels is the subset of cloudflared the routes need.
type Tunnels interface {
	ListTunnels(ctx context.Context) ([]cloudflared.Tunnel, error)
	RouteDNS(ctx context.Context, tunnel, hostname string) error
	ValidateIngress(ctx context.Context, path string) (string, error)
}

// Service runs lifecycle actions on the cloudflared service.
type Service interface {
	Run(ctx context.Context, action service.Action) error
}

// DNS deletes CNAME records when an API token is configured.
type DNS interface {
	Configured() bool
	DeleteCNAME(ctx context.Context, hostname string) (int, error)
}

// Reporter prints progress for the user.
type Reporter interface {
	Printf(format string, args ...any)
}

// Options tunes which side effects run after editing the config.
type Options struct {
	SkipDNS     bool
	SkipRestart bool
}
