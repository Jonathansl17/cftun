package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/store"
)

var errCancelled = errors.New(msg.ErrCancelled)

func newInstallCmd(a *App) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "install",
		Short: msg.InstallShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ensureInstalled(cmd.Context(), a, force)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, msg.FlagForceInstall)
	return cmd
}

func ensureInstalled(ctx context.Context, a *App, force bool) error {
	if version, err := a.Tunnels.Version(ctx); err == nil && !force {
		a.Printf(msg.InfoAlreadyInstalled, version)
		return nil
	}
	family, err := a.Installer.Family()
	if err != nil {
		return err
	}
	a.Printf(msg.InfoInstalling, family)
	if err := a.Installer.Install(ctx); err != nil {
		return err
	}
	version, err := a.Tunnels.Version(ctx)
	if err != nil {
		return err
	}
	a.Printf(msg.InfoInstalled, version)
	return nil
}

func newUninstallCmd(a *App) *cobra.Command {
	var yes, keepCftun bool
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: msg.UninstallShort,
		Long:  msg.UninstallLong,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !yes {
				ok, err := a.Prompt.Confirm(msg.ConfirmUninstall)
				if err != nil {
					return err
				}
				if !ok {
					return errCancelled
				}
			}
			return uninstall(cmd.Context(), a, keepCftun)
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, msg.FlagYes)
	cmd.Flags().BoolVar(&keepCftun, "keep-cftun", false, msg.FlagKeepCftun)
	return cmd
}

func uninstall(ctx context.Context, a *App, keepCftun bool) error {
	params, err := teardownParams(a)
	if err != nil {
		return err
	}
	if err := a.Teardown.Run(ctx, params); err != nil {
		return teardownFailure(err)
	}
	if err := a.Tokens.Clear(); err != nil {
		return err
	}
	a.Printf(msg.InfoUninstalled)
	if keepCftun {
		return nil
	}
	return removeCftun(ctx, a)
}

func errConfigMissing(a *App) error {
	return fmt.Errorf("%s: %w", a.ConfigPath, store.ErrMissing)
}
