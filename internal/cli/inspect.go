package cli

import (
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
)

const (
	tableMinWidth = 0
	tableTabWidth = 4
	tablePadding  = 2
	tablePadChar  = ' '
)

func newListCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   msg.ListShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			rules, err := a.Routes.List()
			if err != nil {
				return err
			}
			if len(rules) == 0 {
				a.Printf(msg.InfoNoRules)
				return nil
			}
			w := tabwriter.NewWriter(a.Out, tableMinWidth, tableTabWidth, tablePadding, tablePadChar, 0)
			fmt.Fprintln(w, msg.ListHeader)
			for _, r := range rules {
				local := a.Health.Probe(cmd.Context(), r.Service)
				fmt.Fprintf(w, msg.ListRowFormat, r.Hostname, r.Service, statusLabel(local))
			}
			return w.Flush()
		},
	}
}

func newCheckCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: msg.CheckShort,
		RunE: func(cmd *cobra.Command, _ []string) error {
			rules, err := a.Routes.List()
			if err != nil {
				return err
			}
			for _, r := range rules {
				checkRule(cmd.Context(), a, r)
			}
			return nil
		},
	}
}

func newValidateCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:         "validate",
		Short:       msg.ValidateShort,
		Annotations: needsCloudflared,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return a.Routes.Validate(cmd.Context())
		},
	}
}

func checkRule(ctx context.Context, a *App, r ingress.Rule) {
	a.Printf(msg.CheckTitle, r.Hostname)
	for _, res := range []health.Result{
		a.Health.Probe(ctx, r.Service),
		a.Health.Probe(ctx, health.PublicURL(r.Hostname)),
	} {
		a.Printf(msg.CheckLine, res.URL, statusLabel(res))
		if hint := res.Hint(); hint != health.HintNone {
			a.Printf(msg.CheckHint, hintTexts[hint])
		}
	}
}

func statusLabel(r health.Result) string {
	switch {
	case r.Err != nil:
		return msg.StatusDown
	case r.OK():
		return fmt.Sprintf(msg.StatusOKFormat, r.Status)
	default:
		return fmt.Sprintf(msg.StatusFailFormat, r.Status)
	}
}
