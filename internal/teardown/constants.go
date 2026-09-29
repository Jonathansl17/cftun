package teardown

import "github.com/Jonathansl17/cftun/internal/paths"

var systemPaths = []string{
	paths.ConfigDir,
	paths.LocalConfigDir,
	paths.RootUserDir,
	paths.LogFile,
	paths.ErrorLogFile,
	paths.ServiceUnit,
	paths.UpdateServiceUnit,
	paths.UpdateTimerUnit,
	paths.InitScript,
}
