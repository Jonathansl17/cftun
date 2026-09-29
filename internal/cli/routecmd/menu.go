package routecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
)

func mark(cmd *cobra.Command, label string, order int) *cobra.Command {
	return uikit.MarkMenu(cmd, uikit.MenuSlot{Group: uikit.GroupRoutes, Label: label, Order: order})
}
