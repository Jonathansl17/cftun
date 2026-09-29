package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
)

const annotationNeedsCloudflared = "cftun/needs-cloudflared"

var needsCloudflared = map[string]string{annotationNeedsCloudflared: "true"}

var errCloudflaredMissing = errors.New(msg.ErrCloudflaredMissing)

func requiresCloudflared(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Annotations[annotationNeedsCloudflared] != "" {
			return true
		}
	}
	return false
}

func ensureRequirements(cmd *cobra.Command, a *App) error {
	if !requiresCloudflared(cmd) {
		return nil
	}
	if _, err := a.Tunnels.Version(cmd.Context()); err == nil {
		return nil
	}
	ok, err := a.Prompt.Confirm(msg.ConfirmInstallCloudflared)
	if err != nil {
		return err
	}
	if !ok {
		return errCloudflaredMissing
	}
	return ensureInstalled(cmd.Context(), a, false)
}
