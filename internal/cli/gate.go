package cli

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g Gate) Require(cmd *cobra.Command) error {
	if !uikit.RequiresCloudflared(cmd) {
		return nil
	}
	if _, err := g.Tunnels.Version(cmd.Context()); err == nil {
		return nil
	}
	ok, err := g.Session.Prompt.Confirm(msg.ConfirmInstallCloudflared)
	if err != nil {
		return err
	}
	if !ok {
		return uikit.ErrCloudflaredMissing
	}
	return g.Install(cmd.Context())
}
