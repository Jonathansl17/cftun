package cli

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cli/setupcmd"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func setupSteps(g groups) []setupcmd.Step {
	return []setupcmd.Step{
		{Title: msg.SetupStepInstall, Run: g.ensureCloudflared},
		{Title: msg.SetupStepLogin, Run: g.auth.EnsureLogin},
		{Title: msg.SetupStepTunnel, Run: g.tunnels.EnsureTunnelConfig},
		{Title: msg.SetupStepService, Run: g.service.InstallService},
		{Title: msg.SetupStepToken, Run: func(context.Context) error { return g.auth.OfferToken() }},
		{Title: msg.SetupStepRoutes, Run: g.routes.OfferRoutes},
	}
}
