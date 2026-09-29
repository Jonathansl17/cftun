package routecmd

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/routes"
)

type RouteService interface {
	List() ([]ingress.Rule, error)
	Add(ctx context.Context, rule ingress.Rule, opts routes.Options) (routes.Outcome, error)
	Remove(ctx context.Context, host string, opts routes.Options) (routes.Outcome, error)
	Edit(ctx context.Context, host string, rule ingress.Rule, opts routes.Options) (routes.Outcome, error)
}

type ConfigValidator interface {
	Validate(ctx context.Context) (string, error)
}

type Prober interface {
	Probe(ctx context.Context, url string) health.Result
}

type Deps struct {
	Session   uikit.Session
	Routes    RouteService
	Validator ConfigValidator
	Prober    Prober
}

type Group struct {
	deps Deps
}

type routeFlags struct {
	Host    string
	NewHost string
	Port    string
	Options routes.Options
}
