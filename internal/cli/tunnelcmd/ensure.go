package tunnelcmd

import (
	"context"
	"errors"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) EnsureTunnelConfig(ctx context.Context) error {
	s := g.deps.Session
	name, err := s.Prompt.Ask(msg.PromptTunnelName, uikit.NotEmpty)
	if err != nil {
		return err
	}
	if err := g.ensureTunnel(ctx, name); err != nil {
		return err
	}
	err = g.initConfig(ctx, name, false)
	if errors.Is(err, uikit.ErrCancelled) {
		s.Printf(msg.InfoConfigKept, g.deps.Config.Path())
		return nil
	}
	return err
}

func (g *Group) ensureTunnel(ctx context.Context, name string) error {
	tunnels, err := g.deps.Tunnels.ListTunnels(ctx)
	if err != nil {
		return err
	}
	if _, ok := cloudflared.FindTunnel(tunnels, name); ok {
		g.deps.Session.Printf(msg.InfoTunnelReused, name)
		return nil
	}
	return g.deps.Tunnels.CreateTunnel(ctx, name)
}
