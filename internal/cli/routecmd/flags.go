package routecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func bindHost(cmd *cobra.Command, flags *routeFlags) {
	cmd.Flags().StringVar(&flags.Host, flagHost, uikit.NoDefault, msg.FlagHost)
}

func bindPort(cmd *cobra.Command, flags *routeFlags) {
	cmd.Flags().StringVar(&flags.Port, flagPort, uikit.NoDefault, msg.FlagPort)
}

func bindOptions(cmd *cobra.Command, flags *routeFlags) {
	cmd.Flags().BoolVar(&flags.Options.SkipDNS, flagNoDNS, false, msg.FlagNoDNS)
	cmd.Flags().BoolVar(&flags.Options.SkipRestart, flagNoRestart, false, msg.FlagNoRestart)
}
