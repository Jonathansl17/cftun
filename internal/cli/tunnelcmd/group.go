package tunnelcmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func New(deps Deps) *Group {
	return &Group{deps: deps}
}

func (g *Group) Commands() []*cobra.Command {
	return []*cobra.Command{g.tunnelsCmd(), g.initCmd()}
}

func (g *Group) tunnelsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: useTunnels, Short: msg.TunnelShort}
	cmd.AddCommand(g.createCmd(), g.listCmd(), g.deleteCmd())
	return uikit.RequireCloudflared(cmd)
}
