package cli

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/paths"
)

func newApp(in infra, g groups) *App {
	app := &App{
		Session:   in.session,
		Files:     in.files,
		Store:     in.store,
		Tunnels:   in.tunnels,
		Installer: in.installer,
		Tokens:    in.tokens,
		DNS:       in.dns,
		Teardown:  in.teardown,
		Quick:     in.quick,
		Groups:    g,
	}
	app.Gate = Gate{Session: in.session, Tunnels: in.tunnels, Install: app.installCloudflared}
	return app
}

func (a *App) installCloudflared(ctx context.Context) error {
	return ensureInstalled(ctx, a, false)
}

func (a *App) userCloudflaredDir() (string, error) {
	home, err := a.Files.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, paths.UserDirName), nil
}

func (a *App) commands() []*cobra.Command {
	return []*cobra.Command{
		newSetupCmd(a), newInstallCmd(a), newUninstallCmd(a), newLoginCmd(a), newTokenCmd(a), newQuickTunnelCmd(a),
	}
}
