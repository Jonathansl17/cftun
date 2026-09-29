package authcmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/dnsapi"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) tokenCmd() *cobra.Command {
	cmd := &cobra.Command{Use: useToken, Short: msg.TokenShort, Long: msg.TokenLong}
	cmd.AddCommand(g.tokenSetCmd(), g.tokenClearCmd(), g.tokenStatusCmd())
	return cmd
}

func (g *Group) tokenClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useTokenClear,
		Short: msg.TokenClearShort,
		RunE: func(*cobra.Command, []string) error {
			return g.deps.Tokens.Clear()
		},
	}
	return mark(cmd, msg.MenuTokenClear, orderTokenClear)
}

func (g *Group) tokenStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useTokenStatus,
		Short: msg.TokenStatusShort,
		RunE:  g.tokenStatus,
	}
	return mark(cmd, msg.MenuTokenStatus, orderTokenStatus)
}

func (g *Group) tokenStatus(*cobra.Command, []string) error {
	token, err := g.deps.Tokens.Load()
	if err != nil {
		return err
	}
	if token != "" {
		g.deps.Session.Printf(msg.InfoTokenSet)
	} else {
		g.deps.Session.Printf(msg.InfoTokenMissing, dnsapi.TokenEnv)
	}
	return nil
}
