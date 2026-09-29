package cli

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/routecmd"
	"github.com/Jonathansl17/cftun/internal/cli/servicecmd"
	"github.com/Jonathansl17/cftun/internal/cli/tunnelcmd"
)

func newGroups(in infra) groups {
	return groups{
		routes: routecmd.New(routecmd.Deps{
			Session: in.session, Routes: in.editor, Validator: in.validator, Prober: in.health,
		}),
		tunnels: tunnelcmd.New(tunnelcmd.Deps{
			Session: in.session, Tunnels: in.tunnels, Initializer: in.initializer, Config: in.store,
		}),
		service: servicecmd.New(servicecmd.Deps{
			Session: in.session, Actions: in.svc, Registry: in.svc, Cloudflared: in.tunnels, Config: in.store,
		}),
	}
}

func (g groups) commands() []*cobra.Command {
	cmds := g.tunnels.Commands()
	cmds = append(cmds, g.routes.Commands()...)
	return append(cmds, g.service.Commands()...)
}
