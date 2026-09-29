package quickcmd

import (
	"context"
	"errors"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/quicktunnel"
)

func (g *Group) Run(ctx context.Context, port string) error {
	s := g.deps.Session
	p, err := uikit.AskPort(s, port, msg.PromptPort)
	if err != nil {
		return err
	}
	if err := g.deps.Ensurer.EnsureInstalled(ctx, noForce); err != nil {
		return err
	}
	service := ingress.LocalService(p)
	s.Printf(msg.InfoQuickStarting, service)
	err = g.deps.Tunnel.Run(ctx, service)
	if errors.Is(err, quicktunnel.ErrInterrupted) {
		s.Printf(msg.InfoQuickStopped)
		return nil
	}
	return err
}
