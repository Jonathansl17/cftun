package authcmd

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/paths"
)

func (g *Group) loginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useLogin,
		Short: msg.LoginShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.EnsureLogin(cmd.Context())
		},
	}
	return mark(uikit.RequireCloudflared(cmd), msg.MenuLogin, orderLogin)
}

func (g *Group) EnsureLogin(ctx context.Context) error {
	dir, err := uikit.UserCloudflaredDir(g.deps.Home.Home)
	if err != nil {
		return err
	}
	cert := filepath.Join(dir, paths.CertFile)
	if g.deps.Files.Exists(cert) {
		g.deps.Session.Printf(msg.InfoAlreadyLoggedIn, cert)
		return nil
	}
	return g.deps.Cloudflared.Login(ctx)
}
