package routecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     useList,
		Aliases: listAliases,
		Short:   msg.ListShort,
		RunE:    g.list,
	}
	return mark(cmd, msg.MenuList, orderList)
}

func (g *Group) list(cmd *cobra.Command, _ []string) error {
	rules, err := g.deps.Routes.List()
	if err != nil {
		return err
	}
	if len(rules) == 0 {
		g.deps.Session.Printf(msg.InfoNoRules)
		return nil
	}
	rows := make([][]string, len(rules))
	for i, r := range rules {
		local := g.deps.Prober.Probe(cmd.Context(), r.Service)
		rows[i] = []string{r.Hostname, r.Service, statusLabel(local)}
	}
	return uikit.PrintTable(g.deps.Session, msg.ListHeader, rows)
}
