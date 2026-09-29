package routecmd

import (
	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/msg"
)

const (
	useAdd      = "add"
	useRemove   = "rm"
	useEdit     = "edit"
	useList     = "list"
	useCheck    = "check"
	useValidate = "validate"

	flagHost      = "host"
	flagNewHost   = "new-host"
	flagPort      = "port"
	flagNoDNS     = "no-dns"
	flagNoRestart = "no-restart"

	noDefault = ""
	noLabel   = ""

	orderList     = 1
	orderAdd      = 2
	orderEdit     = 3
	orderRemove   = 4
	orderCheck    = 5
	orderValidate = 6
)

var (
	removeAliases = []string{"remove", "delete"}
	listAliases   = []string{"ls"}
)

var hintTexts = map[health.Hint]string{
	health.HintConnectionRefused: msg.HintConnectionRefused,
	health.HintDNS:               msg.HintDNS1016,
	health.HintTunnelDown:        msg.HintTunnelDown,
	health.HintGeneric:           msg.HintGeneric,
}
