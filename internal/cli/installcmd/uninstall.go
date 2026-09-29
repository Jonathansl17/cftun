package installcmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) uninstallCmd() *cobra.Command {
	var yes, keepCftun bool
	cmd := &cobra.Command{
		Use:   useUninstall,
		Short: msg.UninstallShort,
		Long:  msg.UninstallLong,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !yes {
				if err := g.deps.Session.ConfirmOrCancel(msg.ConfirmUninstall); err != nil {
					return err
				}
			}
			return g.uninstall(cmd.Context(), keepCftun)
		},
	}
	cmd.Flags().BoolVarP(&yes, flagYes, shortYes, false, msg.FlagYes)
	cmd.Flags().BoolVar(&keepCftun, flagKeepCftun, false, msg.FlagKeepCftun)
	slot := uikit.MenuSlot{Group: uikit.GroupNone, Label: msg.MenuUninstall, Order: uikit.OrderUninstall}
	return uikit.MarkMenu(cmd, slot)
}

func (g *Group) uninstall(ctx context.Context, keepCftun bool) error {
	params, err := g.teardownParams()
	if err != nil {
		return err
	}
	if err := g.deps.Procedure.Run(ctx, params); err != nil {
		return teardownFailure(err)
	}
	if err := g.deps.Tokens.Clear(); err != nil {
		return err
	}
	g.deps.Session.Printf(msg.InfoUninstalled)
	if keepCftun {
		return nil
	}
	return g.removeSelf(ctx)
}
