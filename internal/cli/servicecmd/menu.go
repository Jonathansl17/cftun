package servicecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
)

func mark(cmd *cobra.Command, label string, order int) *cobra.Command {
	return uikit.MarkMenu(cmd, uikit.MenuSlot{Group: uikit.GroupService, Label: label, Order: order})
}
