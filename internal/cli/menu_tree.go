package cli

import (
	"sort"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
)

func buildMenu(root *cobra.Command) ([]menuNode, error) {
	grouped := map[uikit.MenuGroup][]menuNode{}
	if err := collectSlots(root, grouped); err != nil {
		return nil, err
	}
	nodes := grouped[uikit.GroupNone]
	for _, info := range menuGroups {
		if children := grouped[info.group]; len(children) > 0 {
			nodes = append(nodes, menuNode{label: info.label, order: info.order, children: sortNodes(children)})
		}
	}
	return sortNodes(nodes), nil
}

func collectSlots(parent *cobra.Command, grouped map[uikit.MenuGroup][]menuNode) error {
	for _, cmd := range parent.Commands() {
		slot, marked, err := uikit.SlotOf(cmd)
		if err != nil {
			return err
		}
		if marked {
			node := menuNode{label: slot.Label, order: slot.Order, cmd: cmd}
			grouped[slot.Group] = append(grouped[slot.Group], node)
		}
		if err := collectSlots(cmd, grouped); err != nil {
			return err
		}
	}
	return nil
}

func sortNodes(nodes []menuNode) []menuNode {
	sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].order < nodes[j].order })
	return nodes
}
