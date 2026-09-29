package authcmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
)

func New(deps Deps) *Group {
	return &Group{deps: deps}
}

func (g *Group) Commands() []*cobra.Command {
	return []*cobra.Command{g.loginCmd(), g.tokenCmd()}
}

func mark(cmd *cobra.Command, label string, order int) *cobra.Command {
	return uikit.MarkMenu(cmd, uikit.MenuSlot{Group: uikit.GroupAccount, Label: label, Order: order})
}
