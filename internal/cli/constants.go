package cli

import (
	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/paths"
)

var completionFiles = []string{
	paths.BashCompletion,
	paths.ZshCompletion,
	paths.FishCompletion,
}

var hintTexts = map[health.Hint]string{
	health.HintConnectionRefused: msg.HintConnectionRefused,
	health.HintDNS:               msg.HintDNS1016,
	health.HintTunnelDown:        msg.HintTunnelDown,
	health.HintGeneric:           msg.HintGeneric,
}
