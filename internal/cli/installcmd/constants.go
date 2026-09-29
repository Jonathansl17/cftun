package installcmd

import (
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

const (
	useInstall   = "install"
	useUninstall = "uninstall"

	flagForce     = "force"
	flagYes       = "yes"
	shortYes      = "y"
	flagKeepCftun = "keep-cftun"

	orderInstall = 1

	resolveFormat     = "resolve path: %w"
	stepTextFormat    = "%s: %s"
	stepSubjectFormat = "%s %s"
)

var completionFiles = []string{
	paths.BashCompletion,
	paths.ZshCompletion,
	paths.FishCompletion,
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
