package teardown

import "github.com/Jonathansl17/cftun/internal/paths"

const (
	removeBinary   = "rm"
	recursiveFlags = "-rf"
	plainFlags     = "-f"
	endOfOptions   = "--"

	stepErrorFormat    = "step %d: %v"
	unsafeTargetFormat = "unsafe removal target %q"
)

var systemTargets = []Target{
	{Path: paths.ConfigDir, Kind: TargetDirectory},
	{Path: paths.LocalConfigDir, Kind: TargetDirectory},
	{Path: paths.RootUserDir, Kind: TargetDirectory},
	{Path: paths.LogFile, Kind: TargetFile},
	{Path: paths.ErrorLogFile, Kind: TargetFile},
	{Path: paths.ServiceUnit, Kind: TargetFile},
	{Path: paths.UpdateServiceUnit, Kind: TargetFile},
	{Path: paths.UpdateTimerUnit, Kind: TargetFile},
	{Path: paths.InitScript, Kind: TargetFile},
}
