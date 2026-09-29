package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/sysexec"
)

const emptyConfigPattern = "cftun-quick-*.yml"

func newQuickTunnelCmd(a *App) *cobra.Command {
	var port string
	cmd := &cobra.Command{
		Use:   "tunnel",
		Short: msg.QuickShort,
		Long:  msg.QuickLong,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runQuickTunnel(cmd.Context(), a, port)
		},
	}
	cmd.Flags().StringVar(&port, "port", "", msg.FlagPort)
	return cmd
}

// runQuickTunnel exposes a local port on a random trycloudflare.com URL
// until the user presses the interrupt key.
func runQuickTunnel(ctx context.Context, a *App, port string) error {
	p, err := askPort(a, port, msg.PromptPort)
	if err != nil {
		return err
	}
	if err := ensureInstalled(ctx, a, false); err != nil {
		return err
	}
	empty, err := os.CreateTemp("", emptyConfigPattern)
	if err != nil {
		return fmt.Errorf("create empty config: %w", err)
	}
	defer os.Remove(empty.Name())
	if err := empty.Close(); err != nil {
		return fmt.Errorf("close empty config: %w", err)
	}
	service := ingress.LocalService(p)
	a.Printf(msg.InfoQuickStarting, service)
	err = a.Tunnels.QuickTunnel(ctx, service, empty.Name())
	if sysexec.Interrupted(err) {
		a.Printf(msg.InfoQuickStopped)
		return nil
	}
	return err
}
