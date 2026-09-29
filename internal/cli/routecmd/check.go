package routecmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) checkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useCheck,
		Short: msg.CheckShort,
		RunE:  g.check,
	}
	return uikit.MarkIn(uikit.GroupRoutes, cmd, msg.MenuCheck, orderCheck)
}

func (g *Group) check(cmd *cobra.Command, _ []string) error {
	rules, err := g.deps.Routes.List()
	if err != nil {
		return err
	}
	for _, r := range rules {
		g.checkRule(cmd.Context(), r)
	}
	return nil
}

func (g *Group) checkRule(ctx context.Context, r ingress.Rule) {
	s := g.deps.Session
	s.Printf(msg.CheckTitle, r.Hostname)
	for _, res := range []health.Result{
		g.deps.Prober.Probe(ctx, r.Service),
		g.deps.Prober.Probe(ctx, health.PublicURL(r.Hostname)),
	} {
		s.Printf(msg.CheckLine, res.URL, statusLabel(res))
		if hint := res.Hint(); hint != health.HintNone {
			s.Printf(msg.CheckHint, hintTexts[hint])
		}
	}
}
