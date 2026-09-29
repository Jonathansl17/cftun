package servicecmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/apperr"
	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/service"
	"github.com/Jonathansl17/cftun/internal/store"
)

func (g *Group) installCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useInstall,
		Short: msg.ServiceInstallShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.InstallService(cmd.Context())
		},
	}
	return mark(uikit.RequireCloudflared(cmd), msg.MenuServiceInstall, orderInstall)
}

func (g *Group) InstallService(ctx context.Context) error {
	if !g.deps.Config.Exists() {
		return apperr.Wrap(g.deps.Config.Path(), store.ErrMissing)
	}
	if g.deps.Registry.Registered() {
		g.deps.Session.Printf(msg.InfoServiceExists)
	} else if err := g.deps.Cloudflared.InstallService(ctx); err != nil {
		return err
	}
	if err := g.deps.Actions.Run(ctx, service.Enable); err != nil {
		g.deps.Session.Printf(msg.WarnStepFailed, service.Enable, err)
	}
	return g.deps.Actions.Restart(ctx)
}
