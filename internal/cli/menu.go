package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/prompt"
)

func (m menu) run(root *cobra.Command) error {
	nodes, err := buildMenu(root)
	if err != nil {
		return err
	}
	return m.browse(root, menuScreen{title: msg.MenuTitle, nodes: nodes, leave: msg.MenuExit})
}

func (m menu) browse(root *cobra.Command, screen menuScreen) error {
	choices := make([]string, 0, len(screen.nodes)+1)
	for _, node := range screen.nodes {
		choices = append(choices, node.label)
	}
	choices = append(choices, screen.leave)
	for {
		i, err := m.session.Prompt.Select(screen.title, choices)
		if errors.Is(err, prompt.ErrAborted) || (err == nil && i == len(screen.nodes)) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := m.dispatch(root, screen.nodes[i]); err != nil {
			return err
		}
		if m.finished() {
			return nil
		}
	}
}

func (m menu) dispatch(root *cobra.Command, node menuNode) error {
	if node.cmd == nil {
		return m.browse(root, menuScreen{title: node.label, nodes: node.children, leave: msg.MenuBack})
	}
	node.cmd.SetContext(root.Context())
	err := m.gate.Require(node.cmd)
	if err == nil {
		err = node.cmd.RunE(node.cmd, nil)
	}
	if err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			return err
		}
		m.session.PrintError(err)
	}
	return nil
}
