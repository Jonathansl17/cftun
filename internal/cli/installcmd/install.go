package installcmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) installCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   useInstall,
		Short: msg.InstallShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.EnsureInstalled(cmd.Context(), force)
		},
	}
	cmd.Flags().BoolVar(&force, flagForce, false, msg.FlagForceInstall)
	slot := uikit.MenuSlot{Group: uikit.GroupAccount, Label: msg.MenuInstall, Order: uikit.OrderAccountInstall}
	return uikit.MarkMenu(cmd, slot)
}

func (g *Group) EnsureInstalled(ctx context.Context, force bool) error {
	s := g.deps.Session
	if version, err := g.deps.Cloudflared.Version(ctx); err == nil && !force {
		s.Printf(msg.InfoAlreadyInstalled, version)
		return nil
	}
	family, err := g.deps.Installer.Family()
	if err != nil {
		return err
	}
	s.Printf(msg.InfoInstalling, family)
	if err := g.deps.Installer.Install(ctx); err != nil {
		return err
	}
	version, err := g.deps.Cloudflared.Version(ctx)
	if err != nil {
		return err
	}
	s.Printf(msg.InfoInstalled, version)
	return nil
}
