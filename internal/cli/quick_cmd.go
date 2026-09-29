package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/quicktunnel"
)

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

func runQuickTunnel(ctx context.Context, a *App, port string) error {
	p, err := uikit.AskPort(a.Session, port, msg.PromptPort)
	if err != nil {
		return err
	}
	if err := ensureInstalled(ctx, a, false); err != nil {
		return err
	}
	service := ingress.LocalService(p)
	a.Printf(msg.InfoQuickStarting, service)
	err = a.Quick.Run(ctx, service)
	if errors.Is(err, quicktunnel.ErrInterrupted) {
		a.Printf(msg.InfoQuickStopped)
		return nil
	}
	return err
}
