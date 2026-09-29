package cli

import (
	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/routes"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

const (
	resolvePathFormat = "resolve path: %w"
	stepErrorFormat   = "%s: %w"
	stepSubjectFormat = "%s %s"
	subjectTextFormat = "%s: %s"
	causeTextFormat   = "%s: %v"
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

var stepTexts = map[teardown.Step]string{
	teardown.StepStopService:      msg.StepStopService,
	teardown.StepDeleteDNS:        msg.StepDeleteDNS,
	teardown.StepDeleteTunnel:     msg.StepDeleteTunnel,
	teardown.StepUninstallService: msg.StepUninstallService,
	teardown.StepRemovePackage:    msg.StepRemovePackage,
	teardown.StepRemoveFiles:      msg.StepRemoveFiles,
	teardown.StepRemoveTemp:       msg.StepRemoveTemp,
}

var sentinelTexts = map[error]string{
	routes.ErrTunnelNotFound:     msg.ErrTunnelNotFound,
	routes.ErrCredentialsMissing: msg.ErrCredentialsMissing,
}
