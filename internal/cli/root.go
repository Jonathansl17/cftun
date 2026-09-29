package cli

import (
	"io"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
)

const binaryName = "cftun"

func NewRoot(version string, in io.Reader, out io.Writer) *cobra.Command {
	app := &App{}
	var configFlag string
	root := &cobra.Command{
		Use:           binaryName,
		Version:       version,
		Short:         msg.RootShort,
		Long:          msg.RootLong,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			built, err := NewApp(ResolveConfigPath(configFlag), in, out)
			if err != nil {
				return err
			}
			*app = *built
			return ensureRequirements(cmd, app)
		},
		RunE: func(cmd *cobra.Command, _ []string) error { return runMenu(app, cmd.Root()) },
	}
	root.PersistentFlags().StringVar(&configFlag, "config", "", msg.FlagConfig)
	root.SetOut(out)
	root.AddCommand(
		newSetupCmd(app), newInstallCmd(app), newUninstallCmd(app), newLoginCmd(app),
		newTokenCmd(app), newTunnelsCmd(app), newQuickTunnelCmd(app), newInitCmd(app), newAddCmd(app),
		newRemoveCmd(app), newEditCmd(app), newListCmd(app), newCheckCmd(app),
		newValidateCmd(app), newServiceCmd(app),
	)
	return root
}
