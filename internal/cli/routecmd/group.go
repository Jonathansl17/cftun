package routecmd

import "github.com/spf13/cobra"

func New(deps Deps) *Group {
	return &Group{deps: deps}
}

func (g *Group) Commands() []*cobra.Command {
	return []*cobra.Command{
		g.addCmd(), g.removeCmd(), g.editCmd(), g.listCmd(), g.checkCmd(), g.validateCmd(),
	}
}
