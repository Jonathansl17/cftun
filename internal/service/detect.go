package service

import (
	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func Detect(r sysexec.Runner, exists func(string) bool, lookPath func(string) (string, error)) Controller {
	if exists(paths.SystemdMarker) {
		return systemd(r, exists)
	}
	if _, err := lookPath(openRCBinary); err == nil {
		return openRC(r, exists)
	}
	return sysV(r, exists)
}
