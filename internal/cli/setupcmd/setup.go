package setupcmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func New(deps Deps) *Group {
	return &Group{deps: deps}
}

func (g *Group) Commands() []*cobra.Command {
	cmd := &cobra.Command{
		Use:   useSetup,
		Short: msg.SetupShort,
		Long:  msg.SetupLong,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.run(cmd.Context())
		},
	}
	slot := uikit.MenuSlot{Group: uikit.GroupNone, Label: msg.MenuSetup, Order: uikit.OrderSetup}
	return []*cobra.Command{uikit.MarkMenu(cmd, slot)}
}

func (g *Group) run(ctx context.Context) error {
	s := g.deps.Session
	hasDomain, err := s.Prompt.Confirm(msg.ConfirmHasDomain)
	if err != nil {
		return err
	}
	if !hasDomain {
		s.Printf(msg.InfoNoDomain)
		return g.deps.TemporaryTunnel(ctx)
	}
	for i, step := range g.deps.Steps {
		s.Printf(msg.SetupStepFormat, i+1, len(g.deps.Steps), step.Title)
		if err := step.Run(ctx); err != nil {
			return err
		}
	}
	s.Printf(msg.InfoSetupDone)
	return nil
}
