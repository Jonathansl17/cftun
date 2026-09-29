package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/service"
	"github.com/Jonathansl17/cftun/internal/store"
	"github.com/Jonathansl17/cftun/internal/teardown"
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
	a.Printf(msg.InfoInstalling, a.Installer.Family)
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
	var yes bool
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
			return uninstall(cmd.Context(), a)
		},
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, msg.FlagYes)
	return cmd
}

func uninstall(ctx context.Context, a *App) error {
	steps := teardown.Steps{
		Runner:            a.Runner,
		Load:              a.Store.Load,
		DeleteDNS:         a.Routes.DeleteDNS,
		Tunnels:           &a.Tunnels,
		Service:           a.Service,
		Uninstall:         a.Installer.Uninstall,
		Installed:         func(ctx context.Context) bool { _, err := a.Tunnels.Version(ctx); return err == nil },
		ServiceRegistered: func() bool { return service.Installed(a.Files.Exists) },
		Report:            a,
		Paths:             []string{a.ConfigPath, a.Store.BackupPath(), a.UserCloudflaredDir()},
	}
	if err := steps.Run(ctx); err != nil {
		return err
	}
	if err := a.Tokens.Clear(); err != nil {
		return err
	}
	a.Printf(msg.InfoUninstalled)
	return nil
}

func errConfigMissing(a *App) error {
	return fmt.Errorf("%s: %w", a.ConfigPath, store.ErrMissing)
}
