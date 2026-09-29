package routes

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/ingress"
)

type Store interface {
	Exists() bool
	Load() (*ingress.Document, error)
	Save(ctx context.Context, doc *ingress.Document) error
	Restore(ctx context.Context) error
	Delete(ctx context.Context) error
}

type Router interface {
	RouteDNS(ctx context.Context, tunnel, hostname string) error
}

type TunnelCatalog interface {
	ListTunnels(ctx context.Context) ([]cloudflared.Tunnel, error)
}

type IngressValidator interface {
	ValidateIngress(ctx context.Context, path string) (string, error)
}

type ConfigChecker interface {
	Validate(ctx context.Context) (string, error)
}

type PathSource interface {
	Path() string
}

type Restarter interface {
	Restart(ctx context.Context) error
}

type DNSAPI interface {
	Configured() (bool, error)
	DeleteCNAME(ctx context.Context, hostname string) (int, error)
}

type DNSRemover interface {
	Delete(ctx context.Context, host string) (DNSRemoval, error)
}

type HomeLocator interface {
	Home() (string, error)
}

type FileChecker interface {
	Exists(path string) bool
}

type Options struct {
	SkipDNS     bool
	SkipRestart bool
}

type DNSState int

const (
	DNSNotAttempted DNSState = iota
	DNSManual
	DNSDeleted
)

type DNSRemoval struct {
	State   DNSState
	Deleted int
}

type Outcome struct {
	Validation     string
	DNS            DNSRemoval
	RestartSkipped bool
}

type Editor struct {
	Store     Store
	Checker   ConfigChecker
	Router    Router
	DNS       DNSRemover
	Restarter Restarter
}

type Initializer struct {
	Store   Store
	Checker ConfigChecker
	Catalog TunnelCatalog
	Home    HomeLocator
	Files   FileChecker
}

type Validator struct {
	Path    PathSource
	Tunnels IngressValidator
}

type DNSCleaner struct {
	API DNSAPI
}

type RollbackError struct {
	Cause error
}

type DNSRollbackError struct {
	Cause error
}

type applier struct {
	store   Store
	checker ConfigChecker
}
