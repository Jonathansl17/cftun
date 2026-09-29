package cli

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/dnsapi"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/paths"
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
	return filepath.Join(a.UserCloudflaredDir(), paths.CertFile)
}

func loggedIn(a *App) bool {
	return a.Files.Exists(certPath(a))
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
	token, err := valueOrAsk(a, value, msg.PromptToken, notEmpty)
	if err != nil {
		return err
	}
	if err := a.Tokens.Save(token); err != nil {
		return err
	}
	a.Printf(msg.InfoTokenSaved, a.Tokens.Location())
	return nil
}
