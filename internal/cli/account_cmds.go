package cli

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/dnsapi"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/paths"
)

func newLoginCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: msg.LoginShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return ensureLogin(cmd.Context(), a)
		},
	}
	return uikit.RequireCloudflared(cmd)
}

func ensureLogin(ctx context.Context, a *App) error {
	dir, err := a.userCloudflaredDir()
	if err != nil {
		return err
	}
	cert := filepath.Join(dir, paths.CertFile)
	if a.Files.Exists(cert) {
		a.Printf(msg.InfoAlreadyLoggedIn, cert)
		return nil
	}
	return a.Tunnels.Login(ctx)
}

func newTokenCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{Use: "token", Short: msg.TokenShort, Long: msg.TokenLong}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "set [token]",
			Short: msg.TokenSetShort,
			Args:  cobra.MaximumNArgs(1),
			RunE:  func(_ *cobra.Command, args []string) error { return setToken(a, uikit.FirstArg(args)) },
		},
		&cobra.Command{
			Use:   "clear",
			Short: msg.TokenClearShort,
			RunE: func(*cobra.Command, []string) error {
				return a.Tokens.Clear()
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: msg.TokenStatusShort,
			RunE: func(*cobra.Command, []string) error {
				configured, err := a.DNS.Configured()
				if err != nil {
					return err
				}
				if configured {
					a.Printf(msg.InfoTokenSet)
				} else {
					a.Printf(msg.InfoTokenMissing, dnsapi.TokenEnv)
				}
				return nil
			},
		},
	)
	return cmd
}

func setToken(a *App, value string) error {
	token, err := a.ValueOrAsk(value, msg.PromptToken, uikit.NotEmpty)
	if err != nil {
		return err
	}
	if err := a.Tokens.Save(token); err != nil {
		return err
	}
	a.Printf(msg.InfoTokenSaved, a.Tokens.Location())
	return nil
}
