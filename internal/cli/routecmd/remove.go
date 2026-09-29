package routecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) removeCmd() *cobra.Command {
	var flags routeFlags
	cmd := &cobra.Command{
		Use:     useRemove,
		Aliases: removeAliases,
		Short:   msg.RemoveShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.removeRoute(cmd, flags)
		},
	}
	bindHost(cmd, &flags)
	bindOptions(cmd, &flags)
	return uikit.RequireCloudflared(cmd)
}

func (g *Group) removeRoute(cmd *cobra.Command, flags routeFlags) error {
	rule, err := g.pickRule(flags.Host)
	if err != nil {
		return err
	}
	out, err := g.deps.Routes.Remove(cmd.Context(), rule.Hostname, flags.Options)
	g.printOutcome(rule.Hostname, out)
	if err != nil {
		return err
	}
	g.deps.Session.Printf(msg.InfoRemoved, rule.Hostname)
	return nil
}
