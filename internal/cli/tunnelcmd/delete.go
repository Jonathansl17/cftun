package tunnelcmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useDelete,
		Short: msg.TunnelDeleteShort,
		Args:  cobra.MaximumNArgs(maxNameArgs),
		RunE:  g.delete,
	}
	return uikit.MarkIn(uikit.GroupTunnels, cmd, msg.MenuTunnelDelete, orderDelete)
}

func (g *Group) delete(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	name := uikit.FirstArg(args)
	if name == "" {
		picked, found, err := g.pickTunnel(ctx)
		if err != nil || !found {
			return err
		}
		name = picked
	}
	if err := g.deps.Session.ConfirmOrCancel(fmt.Sprintf(msg.ConfirmDeleteTunnel, name)); err != nil {
		return err
	}
	return g.deps.Tunnels.DeleteTunnel(ctx, name)
}

func (g *Group) pickTunnel(ctx context.Context) (string, bool, error) {
	tunnels, err := g.deps.Tunnels.ListTunnels(ctx)
	if err != nil {
		return "", false, err
	}
	if len(tunnels) == 0 {
		g.deps.Session.Printf(msg.InfoNoTunnels)
		return "", false, nil
	}
	labels := make([]string, len(tunnels))
	for i, t := range tunnels {
		labels[i] = t.Name
	}
	i, err := g.deps.Session.Prompt.Select(msg.PromptPickTunnel, labels)
	if err != nil {
		return "", false, err
	}
	return tunnels[i].Name, true, nil
}
