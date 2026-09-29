package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
)

func newInitCmd(a *App) *cobra.Command {
	var tunnel string
	var force bool
	cmd := &cobra.Command{
		Use:         "init",
		Short:       msg.InitShort,
		Annotations: needsCloudflared,
		RunE: func(cmd *cobra.Command, _ []string) error {
			name, err := valueOrAsk(a, tunnel, msg.PromptTunnelName, notEmpty)
			if err != nil {
				return err
			}
			return initConfig(cmd.Context(), a, name, force)
		},
	}
	cmd.Flags().StringVar(&tunnel, "tunnel", "", msg.FlagTunnel)
	cmd.Flags().BoolVar(&force, "force", false, msg.FlagForceInit)
	return cmd
}

func initConfig(ctx context.Context, a *App, tunnel string, force bool) error {
	if a.Store.Exists() && !force {
		ok, err := a.Prompt.Confirm(msg.ConfirmOverwrite)
		if err != nil || !ok {
			return orCancelled(err)
		}
	}
	validation, err := a.Initializer.Init(ctx, tunnel)
	if err != nil {
		return present(err)
	}
	a.Printf(msg.InfoValidated, validation)
	a.Printf(msg.InfoConfigWritten, a.ConfigPath)
	return nil
}
