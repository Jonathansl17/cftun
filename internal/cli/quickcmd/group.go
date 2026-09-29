package quickcmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func New(deps Deps) *Group {
	return &Group{deps: deps}
}

func (g *Group) Commands() []*cobra.Command {
	var port string
	cmd := &cobra.Command{
		Use:   useTunnel,
		Short: msg.QuickShort,
		Long:  msg.QuickLong,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return g.Run(cmd.Context(), port)
		},
	}
	cmd.Flags().StringVar(&port, flagPort, uikit.NoDefault, msg.FlagPort)
	slot := uikit.MenuSlot{Group: uikit.GroupNone, Label: msg.MenuQuick, Order: uikit.OrderQuick}
	return []*cobra.Command{uikit.MarkMenu(cmd, slot)}
}
