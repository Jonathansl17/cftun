package tunnelcmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useList,
		Short: msg.TunnelListShort,
		RunE:  g.list,
	}
	return uikit.MarkIn(uikit.GroupTunnels, cmd, msg.MenuTunnelList, orderList)
}

func (g *Group) list(cmd *cobra.Command, _ []string) error {
	tunnels, err := g.deps.Tunnels.ListTunnels(cmd.Context())
	if err != nil {
		return err
	}
	if len(tunnels) == 0 {
		g.deps.Session.Printf(msg.InfoNoTunnels)
		return nil
	}
	rows := make([][]string, len(tunnels))
	for i, t := range tunnels {
		rows[i] = []string{t.Name, t.ID}
	}
	return uikit.PrintTable(g.deps.Session, msg.TunnelHeader, rows)
}
