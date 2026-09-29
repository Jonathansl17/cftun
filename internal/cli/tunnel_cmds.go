package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
)

func newTunnelCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{Use: "tunnel", Short: msg.TunnelShort}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "create [name]",
			Short: msg.TunnelCreateShort,
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				name, err := valueOrAsk(a, firstArg(args), msg.PromptTunnelName, notEmpty)
				if err != nil {
					return err
				}
				return a.Tunnels.CreateTunnel(cmd.Context(), name)
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: msg.TunnelListShort,
			RunE:  func(cmd *cobra.Command, _ []string) error { return listTunnels(cmd, a) },
		},
		&cobra.Command{
			Use:   "delete [name]",
			Short: msg.TunnelDeleteShort,
			Args:  cobra.MaximumNArgs(1),
			RunE:  func(cmd *cobra.Command, args []string) error { return deleteTunnel(cmd, a, firstArg(args)) },
		},
	)
	return cmd
}

func listTunnels(cmd *cobra.Command, a *App) error {
	tunnels, err := a.Tunnels.ListTunnels(cmd.Context())
	if err != nil {
		return err
	}
	if len(tunnels) == 0 {
		a.Printf(msg.InfoNoTunnels)
		return nil
	}
	w := tabwriter.NewWriter(a.Out, tableMinWidth, tableTabWidth, tablePadding, tablePadChar, 0)
	fmt.Fprintln(w, msg.TunnelHeader)
	for _, t := range tunnels {
		fmt.Fprintf(w, msg.TunnelRowFormat, t.Name, t.ID)
	}
	return w.Flush()
}

func deleteTunnel(cmd *cobra.Command, a *App, name string) error {
	if name == "" {
		tunnels, err := a.Tunnels.ListTunnels(cmd.Context())
		if err != nil {
			return err
		}
		if len(tunnels) == 0 {
			a.Printf(msg.InfoNoTunnels)
			return nil
		}
		labels := make([]string, len(tunnels))
		for i, t := range tunnels {
			labels[i] = t.Name
		}
		i, err := a.Prompt.Select(msg.PromptPickTunnel, labels)
		if err != nil {
			return err
		}
		name = tunnels[i].Name
	}
	ok, err := a.Prompt.Confirm(fmt.Sprintf(msg.ConfirmDeleteTunnel, name))
	if err != nil || !ok {
		return orCancelled(err)
	}
	return a.Tunnels.DeleteTunnel(cmd.Context(), name)
}

// orCancelled returns err, or errCancelled when the user simply said no.
func orCancelled(err error) error {
	if err != nil {
		return err
	}
	return errCancelled
}
