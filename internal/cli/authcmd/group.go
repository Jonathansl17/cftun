package authcmd

import (
	"github.com/spf13/cobra"
)

func New(deps Deps) *Group {
	return &Group{deps: deps}
}

func (g *Group) Commands() []*cobra.Command {
	return []*cobra.Command{g.loginCmd(), g.tokenCmd()}
}
