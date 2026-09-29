package routecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
)

func bindHost(cmd *cobra.Command, flags *routeFlags) {
	cmd.Flags().StringVar(&flags.Host, flagHost, noDefault, msg.FlagHost)
}

func bindPort(cmd *cobra.Command, flags *routeFlags) {
	cmd.Flags().StringVar(&flags.Port, flagPort, noDefault, msg.FlagPort)
}

func bindOptions(cmd *cobra.Command, flags *routeFlags) {
	cmd.Flags().BoolVar(&flags.Options.SkipDNS, flagNoDNS, false, msg.FlagNoDNS)
	cmd.Flags().BoolVar(&flags.Options.SkipRestart, flagNoRestart, false, msg.FlagNoRestart)
}
