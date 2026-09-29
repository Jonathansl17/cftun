package authcmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/dnsapi"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) tokenSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useTokenSet,
		Short: msg.TokenSetShort,
		Args:  cobra.NoArgs,
		RunE:  func(*cobra.Command, []string) error { return g.saveToken() },
	}
	return mark(cmd, msg.MenuTokenSet, orderTokenSet)
}

func (g *Group) saveToken() error {
	s := g.deps.Session
	token, err := s.ValueOrAsk(g.deps.Env(dnsapi.TokenEnv), msg.PromptToken, uikit.NotEmpty)
	if err != nil {
		return err
	}
	if err := g.deps.Tokens.Save(token); err != nil {
		return err
	}
	s.Printf(msg.InfoTokenSaved, g.deps.Tokens.Location())
	return nil
}

func (g *Group) OfferToken() error {
	token, err := g.deps.Tokens.Load()
	if err != nil {
		return err
	}
	if token != "" {
		g.deps.Session.Printf(msg.InfoTokenSet)
		return nil
	}
	g.deps.Session.Printf(msg.TokenLong)
	ok, err := g.deps.Session.Prompt.Confirm(msg.ConfirmSetToken)
	if err != nil || !ok {
		return err
	}
	return g.saveToken()
}
