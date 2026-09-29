package cli

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/service"
)

func newServiceCmd(a *App) *cobra.Command {
	cmd := &cobra.Command{Use: "service", Short: msg.ServiceShort}
	cmd.AddCommand(&cobra.Command{
		Use:         "install",
		Short:       msg.ServiceInstallShort,
		Annotations: needsCloudflared,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return installService(cmd, a)
		},
	})
	for _, action := range service.Actions {
		action := action
		cmd.AddCommand(&cobra.Command{
			Use:   string(action),
			Short: msg.ServiceActionShort[action],
			RunE: func(cmd *cobra.Command, _ []string) error {
				return a.Service.Run(cmd.Context(), action)
			},
		})
	}
	return cmd
}

func installService(cmd *cobra.Command, a *App) error {
	if !a.Store.Exists() {
		return errConfigMissing(a)
	}
	if a.Service.Registered() {
		a.Printf(msg.InfoServiceExists)
	} else if err := a.Tunnels.InstallService(cmd.Context()); err != nil {
		return err
	}
	if err := a.Service.Run(cmd.Context(), service.Enable); err != nil {
		a.Printf(msg.WarnStepFailed, service.Enable, err)
	}
	return a.Service.Restart(cmd.Context())
}
