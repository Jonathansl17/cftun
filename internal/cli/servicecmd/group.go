package servicecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func New(deps Deps) *Group {
	return &Group{deps: deps}
}

func (g *Group) Commands() []*cobra.Command {
	cmd := uikit.RequireCloudflared(&cobra.Command{Use: useService, Short: msg.ServiceShort})
	cmd.AddCommand(g.installCmd())
	cmd.AddCommand(g.actionCmds()...)
	return []*cobra.Command{cmd}
}
