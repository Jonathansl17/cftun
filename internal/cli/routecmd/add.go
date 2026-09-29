package routecmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) addCmd() *cobra.Command {
	var flags routeFlags
	cmd := &cobra.Command{
		Use:   useAdd,
		Short: msg.AddShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.addRoute(cmd.Context(), flags)
		},
	}
	bindHost(cmd, &flags)
	bindPort(cmd, &flags)
	bindOptions(cmd, &flags)
	return uikit.RequireCloudflared(cmd)
}

func (g *Group) addRoute(ctx context.Context, flags routeFlags) error {
	s := g.deps.Session
	hostname, err := uikit.AskHostname(s, flags.Host, msg.PromptHostname)
	if err != nil {
		return err
	}
	port, err := uikit.AskPort(s, flags.Port, msg.PromptPort)
	if err != nil {
		return err
	}
	rule := ingress.Rule{Hostname: hostname, Service: ingress.LocalService(port)}
	out, err := g.deps.Routes.Add(ctx, rule, flags.Options)
	g.printOutcome(rule.Hostname, out)
	if err != nil {
		return err
	}
	s.Printf(msg.InfoAdded, rule.Hostname, rule.Service)
	return nil
}

func (g *Group) OfferRoutes(ctx context.Context) error {
	for {
		ok, err := g.deps.Session.Prompt.Confirm(msg.ConfirmAddRoute)
		if err != nil || !ok {
			return err
		}
		if err := g.addRoute(ctx, routeFlags{}); err != nil {
			g.deps.Session.PrintError(err)
		}
	}
}
