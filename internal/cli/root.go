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
	gate := Gate{Session: in.session, Tunnels: in.tunnels, Install: g.ensureCloudflared}
	m := menu{session: in.session, gate: gate, finished: g.install.SelfRemoved}
	root := &cobra.Command{
		Use:           binaryName,
		Version:       version,
		Short:         msg.RootShort,
		Long:          msg.RootLong,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return gate.Require(cmd)
		},
		RunE: func(cmd *cobra.Command, _ []string) error { return m.run(cmd.Root()) },
	}
	root.PersistentFlags().StringVar(&locator.Flag, flagConfig, noDefault, msg.FlagConfig)
	root.SetOut(streams.Out)
	root.AddCommand(g.commands()...)
	return root, nil
}
