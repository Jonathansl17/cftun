package service

import (
	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func systemd(r sysexec.Runner, exists func(string) bool) Controller {
	return Controller{Name: systemdName, Runner: r, exists: exists, commands: map[Action][]string{
		Enable:  {systemctlBinary, string(Enable), unitName},
		Disable: {systemctlBinary, string(Disable), unitName},
		Start:   {systemctlBinary, string(Start), unitName},
		Stop:    {systemctlBinary, string(Stop), unitName},
		Restart: {systemctlBinary, string(Restart), unitName},
		Status:  {systemctlBinary, string(Status), systemdNoPagerFlag, unitName},
		Logs:    {journalctlBinary, journalUnitFlag, unitName, followLogsFlag},
	}}
}

func openRC(r sysexec.Runner, exists func(string) bool) Controller {
	return Controller{Name: openRCName, Runner: r, exists: exists, commands: map[Action][]string{
		Enable:  {openRCUpdateBinary, openRCAddVerb, unitName, openRCRunlevel},
		Disable: {openRCUpdateBinary, openRCDeleteVerb, unitName, openRCRunlevel},
		Start:   {openRCBinary, unitName, string(Start)},
		Stop:    {openRCBinary, unitName, string(Stop)},
		Restart: {openRCBinary, unitName, string(Restart)},
		Status:  {openRCBinary, unitName, string(Status)},
		Logs:    {tailBinary, followLogsFlag, paths.LogFile},
	}}
}

func sysV(r sysexec.Runner, exists func(string) bool) Controller {
	return Controller{Name: sysVName, Runner: r, exists: exists, commands: map[Action][]string{
		Start:   {sysVBinary, unitName, string(Start)},
		Stop:    {sysVBinary, unitName, string(Stop)},
		Restart: {sysVBinary, unitName, string(Restart)},
		Status:  {sysVBinary, unitName, string(Status)},
		Logs:    {tailBinary, followLogsFlag, paths.LogFile},
	}}
}
