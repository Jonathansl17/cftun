package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/routes"
)

// bindRouteOptions registers the flags shared by add, rm and edit.
func bindRouteOptions(cmd *cobra.Command, opts *routes.Options) {
	cmd.Flags().BoolVar(&opts.SkipDNS, "no-dns", false, msg.FlagNoDNS)
	cmd.Flags().BoolVar(&opts.SkipRestart, "no-restart", false, msg.FlagNoRestart)
}

func newAddCmd(a *App) *cobra.Command {
	var host, port string
	var opts routes.Options
	cmd := &cobra.Command{
		Use:         "add",
		Short:       msg.AddShort,
		Annotations: needsCloudflared,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return addRoute(cmd.Context(), a, host, port, opts)
		},
	}
	cmd.Flags().StringVar(&host, "host", "", msg.FlagHost)
	cmd.Flags().StringVar(&port, "port", "", msg.FlagPort)
	bindRouteOptions(cmd, &opts)
	return cmd
}

func newRemoveCmd(a *App) *cobra.Command {
	var host string
	var opts routes.Options
	cmd := &cobra.Command{
		Use:         "rm",
		Aliases:     []string{"remove", "delete"},
		Short:       msg.RemoveShort,
		Annotations: needsCloudflared,
		RunE: func(cmd *cobra.Command, _ []string) error {
			rule, err := pickRule(a, host)
			if err != nil {
				return err
			}
			if err := a.Routes.Remove(cmd.Context(), rule.Hostname, opts); err != nil {
				return err
			}
			a.Printf(msg.InfoRemoved, rule.Hostname)
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "", msg.FlagHost)
	bindRouteOptions(cmd, &opts)
	return cmd
}

// addRoute asks for any missing value and publishes the route.
func addRoute(ctx context.Context, a *App, host, port string, opts routes.Options) error {
	hostname, err := askHostname(a, host, msg.PromptHostname)
	if err != nil {
		return err
	}
	p, err := askPort(a, port, msg.PromptPort)
	if err != nil {
		return err
	}
	rule := ingress.Rule{Hostname: hostname, Service: ingress.LocalService(p)}
	if err := a.Routes.Add(ctx, rule, opts); err != nil {
		return err
	}
	a.Printf(msg.InfoAdded, rule.Hostname, rule.Service)
	return nil
}
