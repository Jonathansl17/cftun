package tunnelcmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) initCmd() *cobra.Command {
	var tunnel string
	var force bool
	cmd := &cobra.Command{
		Use:   useInit,
		Short: msg.InitShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			name, err := g.deps.Session.ValueOrAsk(tunnel, msg.PromptTunnelName, uikit.NotEmpty)
			if err != nil {
				return err
			}
			return g.initConfig(cmd.Context(), name, force)
		},
	}
	cmd.Flags().StringVar(&tunnel, flagTunnel, noDefault, msg.FlagTunnel)
	cmd.Flags().BoolVar(&force, flagForce, false, msg.FlagForceInit)
	return uikit.RequireCloudflared(cmd)
}

func (g *Group) initConfig(ctx context.Context, tunnel string, force bool) error {
	if g.deps.Config.Exists() && !force {
		if err := g.deps.Session.ConfirmOrCancel(msg.ConfirmOverwrite); err != nil {
			return err
		}
	}
	validation, err := g.deps.Initializer.Init(ctx, tunnel)
	if err != nil {
		return err
	}
	g.deps.Session.Printf(msg.InfoValidated, validation)
	g.deps.Session.Printf(msg.InfoConfigWritten, g.deps.Config.Path())
	return nil
}
