package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/authcmd"
	"github.com/Jonathansl17/cftun/internal/cli/installcmd"
	"github.com/Jonathansl17/cftun/internal/cli/quickcmd"
	"github.com/Jonathansl17/cftun/internal/cli/routecmd"
	"github.com/Jonathansl17/cftun/internal/cli/servicecmd"
	"github.com/Jonathansl17/cftun/internal/cli/setupcmd"
	"github.com/Jonathansl17/cftun/internal/cli/tunnelcmd"
)

func newGroups(in infra) groups {
	g := groups{}
	g.install = installcmd.New(installcmd.Deps{
		Session: in.session, Installer: in.installer, Cloudflared: in.tunnels, Procedure: in.teardown,
		Tokens: in.tokens, Config: in.store, Home: in.files, Files: in.files, Remover: in.store.Writer,
	})
	g.auth = authcmd.New(authcmd.Deps{
		Session: in.session, Cloudflared: in.tunnels, Files: in.files, Home: in.files, Tokens: in.tokens,
		Env: os.Getenv,
	})
	g.quick = quickcmd.New(quickcmd.Deps{Session: in.session, Ensurer: g.install, Tunnel: in.quick})
	g.routes = routecmd.New(routecmd.Deps{
		Session: in.session, Routes: in.editor, Validator: in.validator, Prober: in.health,
	})
	g.tunnels = tunnelcmd.New(tunnelcmd.Deps{
		Session: in.session, Tunnels: in.tunnels, Initializer: in.initializer, Config: in.store,
	})
	g.service = servicecmd.New(servicecmd.Deps{
		Session: in.session, Actions: in.svc, Registry: in.svc, Cloudflared: in.tunnels, Config: in.store,
	})
	g.setup = setupcmd.New(setupcmd.Deps{
		Session: in.session, Steps: setupSteps(g), TemporaryTunnel: g.temporaryTunnel,
	})
	return g
}

func (g groups) commands() []*cobra.Command {
	providers := []commandProvider{g.setup, g.install, g.auth, g.quick, g.tunnels, g.routes, g.service}
	var cmds []*cobra.Command
	for _, provider := range providers {
		cmds = append(cmds, provider.Commands()...)
	}
	return cmds
}

func (g groups) ensureCloudflared(ctx context.Context) error {
	return g.install.EnsureInstalled(ctx, false)
}

func (g groups) temporaryTunnel(ctx context.Context) error {
	return g.quick.Run(ctx, noPort)
}
