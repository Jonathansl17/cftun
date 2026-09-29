package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/routes"
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
		{msg.SetupStepTunnel, func() error { return ensureTunnelConfig(ctx, a) }},
		{msg.SetupStepService, func() error { return installService(cmd, a) }},
		{msg.SetupStepToken, func() error { return offerToken(a) }},
		{msg.SetupStepRoutes, func() error { return offerRoutes(ctx, a) }},
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

func ensureLogin(ctx context.Context, a *App) error {
	if loggedIn(a) {
		a.Printf(msg.InfoAlreadyLoggedIn, certPath(a))
		return nil
	}
	return a.Tunnels.Login(ctx)
}

func ensureTunnelConfig(ctx context.Context, a *App) error {
	name, err := a.Prompt.Ask(msg.PromptTunnelName, notEmpty)
	if err != nil {
		return err
	}
	tunnels, err := a.Tunnels.ListTunnels(ctx)
	if err != nil {
		return err
	}
	if _, ok := cloudflared.FindTunnel(tunnels, name); ok {
		a.Printf(msg.InfoTunnelReused, name)
	} else if err := a.Tunnels.CreateTunnel(ctx, name); err != nil {
		return err
	}
	err = initConfig(ctx, a, name, false)
	if errors.Is(err, errCancelled) {
		a.Printf(msg.InfoConfigKept, a.ConfigPath)
		return nil
	}
	return err
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

func offerRoutes(ctx context.Context, a *App) error {
	for {
		ok, err := a.Prompt.Confirm(msg.ConfirmAddRoute)
		if err != nil || !ok {
			return err
		}
		if err := addRoute(ctx, a, "", "", routes.Options{}); err != nil {
			a.Printf(msg.ErrorFormat, err)
		}
	}
}
