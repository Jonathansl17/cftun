package tunnelcmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useCreate,
		Short: msg.TunnelCreateShort,
		Args:  cobra.MaximumNArgs(maxNameArgs),
		RunE:  g.create,
	}
	return mark(cmd, msg.MenuTunnelCreate, orderCreate)
}

func (g *Group) create(cmd *cobra.Command, args []string) error {
	name, err := g.deps.Session.ValueOrAsk(uikit.FirstArg(args), msg.PromptTunnelName, uikit.NotEmpty)
	if err != nil {
		return err
	}
	return g.deps.Tunnels.CreateTunnel(cmd.Context(), name)
}
