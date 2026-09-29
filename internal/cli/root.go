package cli

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
)

func NewRoot(version string, streams Streams) (*cobra.Command, error) {
	locator := newConfigLocator()
	in, err := newInfra(streams, locator)
	if err != nil {
		return nil, err
	}
	g := newGroups(in)
	app := newApp(in, g)
	root := &cobra.Command{
		Use:           binaryName,
		Version:       version,
		Short:         msg.RootShort,
		Long:          msg.RootLong,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return app.Gate.Require(cmd)
		},
		RunE: func(cmd *cobra.Command, _ []string) error { return runMenu(app, cmd.Root()) },
	}
	root.PersistentFlags().StringVar(&locator.Flag, flagConfig, noDefault, msg.FlagConfig)
	root.SetOut(streams.Out)
	root.AddCommand(app.commands()...)
	root.AddCommand(g.commands()...)
	return root, nil
}
