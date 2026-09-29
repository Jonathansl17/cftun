package routecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) editCmd() *cobra.Command {
	var flags routeFlags
	cmd := &cobra.Command{
		Use:   useEdit,
		Short: msg.EditShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.editRoute(cmd, flags)
		},
	}
	bindHost(cmd, &flags)
	cmd.Flags().StringVar(&flags.NewHost, flagNewHost, noDefault, msg.FlagNewHost)
	bindPort(cmd, &flags)
	bindOptions(cmd, &flags)
	return mark(uikit.RequireCloudflared(cmd), msg.MenuEdit, orderEdit)
}

func (g *Group) editRoute(cmd *cobra.Command, flags routeFlags) error {
	rule, err := g.pickRule(flags.Host)
	if err != nil {
		return err
	}
	updated, err := g.askEdit(rule, flags)
	if err != nil {
		return err
	}
	out, err := g.deps.Routes.Edit(cmd.Context(), rule.Hostname, updated, flags.Options)
	g.printOutcome(rule.Hostname, out)
	if err != nil {
		return err
	}
	g.deps.Session.Printf(msg.InfoAdded, updated.Hostname, updated.Service)
	return nil
}
