package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/dnsapi"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func newLoginCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:         "login",
		Short:       msg.LoginShort,
		Annotations: needsCloudflared,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if loggedIn(a) {
				a.Printf(msg.InfoAlreadyLoggedIn, certPath(a))
				return nil
			}
			return a.Tunnels.Login(cmd.Context())
		},
	}
}

func certPath(a *App) string {
	return filepath.Join(a.UserCloudflaredDir(), cloudflared.CertFile)
}

func loggedIn(a *App) bool {
	_, err := os.Stat(certPath(a))
	return err == nil
}

func newTokenCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{Use: "token", Short: msg.TokenShort, Long: msg.TokenLong}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "set [token]",
			Short: msg.TokenSetShort,
			Args:  cobra.MaximumNArgs(1),
			RunE:  func(_ *cobra.Command, args []string) error { return setToken(a, firstArg(args)) },
		},
		&cobra.Command{
			Use:   "clear",
			Short: msg.TokenClearShort,
			RunE: func(*cobra.Command, []string) error {
				a.setToken("")
				return a.Tokens.Clear()
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: msg.TokenStatusShort,
			RunE: func(*cobra.Command, []string) error {
				if a.DNS.Configured() {
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
	token, err := valueOrAsk(a, value, msg.PromptToken, notEmpty)
	if err != nil {
		return err
	}
	if err := a.Tokens.Save(token); err != nil {
		return err
	}
	a.setToken(token)
	a.Printf(msg.InfoTokenSaved, a.Tokens.Path)
	return nil
}

func (a *App) setToken(token string) {
	a.DNS.Token = token
	a.Routes.DNS = a.DNS
}
