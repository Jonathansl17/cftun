package installcmd

import "github.com/spf13/cobra"

func New(deps Deps) *Group {
	return &Group{deps: deps}
}

func (g *Group) Commands() []*cobra.Command {
	return []*cobra.Command{g.installCmd(), g.uninstallCmd()}
}

func (g *Group) SelfRemoved() bool {
	return g.removed
}
