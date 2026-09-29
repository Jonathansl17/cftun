package cli

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
)

func newSetupCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: msg.SetupShort,
		Long:  msg.SetupLong,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSetup(cmd, a)
		},
	}
}

func runSetup(cmd *cobra.Command, a *App) error {
	ctx := cmd.Context()
	hasDomain, err := a.Prompt.Confirm(msg.ConfirmHasDomain)
	if err != nil {
		return err
	}
	if !hasDomain {
		a.Printf(msg.InfoNoDomain)
		return runQuickTunnel(ctx, a, "")
	}
	steps := []struct {
		title string
		run   func() error
	}{
		{msg.SetupStepInstall, func() error { return ensureInstalled(ctx, a, false) }},
		{msg.SetupStepLogin, func() error { return ensureLogin(ctx, a) }},
		{msg.SetupStepTunnel, func() error { return a.Groups.tunnels.EnsureTunnelConfig(ctx) }},
		{msg.SetupStepService, func() error { return a.Groups.service.InstallService(ctx) }},
		{msg.SetupStepToken, func() error { return offerToken(a) }},
		{msg.SetupStepRoutes, func() error { return a.Groups.routes.OfferRoutes(ctx) }},
	}
	for i, step := range steps {
		a.Printf(msg.SetupStepFormat, i+1, len(steps), step.title)
		if err := step.run(); err != nil {
			return err
		}
	}
	a.Printf(msg.InfoSetupDone)
	return nil
}

func offerToken(a *App) error {
	configured, err := a.DNS.Configured()
	if err != nil {
		return err
	}
	if configured {
		a.Printf(msg.InfoTokenSet)
		return nil
	}
	a.Printf(msg.TokenLong)
	ok, err := a.Prompt.Confirm(msg.ConfirmSetToken)
	if err != nil || !ok {
		return err
	}
	return setToken(a, "")
}
