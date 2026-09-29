package servicecmd

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/service"
)

type Actions interface {
	Run(ctx context.Context, action service.Action) error
	Restart(ctx context.Context) error
}

type Registry interface {
	Registered() bool
}

type Cloudflared interface {
	InstallService(ctx context.Context) error
}

type Config interface {
	Path() string
	Exists() bool
}

type Deps struct {
	Session     uikit.Session
	Actions     Actions
	Registry    Registry
	Cloudflared Cloudflared
	Config      Config
}

type Group struct {
	deps Deps
}
