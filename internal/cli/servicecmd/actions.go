package servicecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/service"
)

func (g *Group) actionCmds() []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(service.Actions))
	for _, action := range service.Actions {
		cmds = append(cmds, g.actionCmd(action))
	}
	return cmds
}

func (g *Group) actionCmd(action service.Action) *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(action),
		Short: msg.ServiceActionShort[action],
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.deps.Actions.Run(cmd.Context(), action)
		},
	}
	return mark(cmd, msg.MenuServiceLabel[action], actionOrder[action])
}
