// Package service controls the cloudflared system service on any init system.
package service

import (
	"context"
	"fmt"

	"github.com/Jonathansl17/cloudfare/internal/sysexec"
)

// Action is one lifecycle operation on the service.
type Action string

// Supported actions.
const (
	Enable  Action = "enable"
	Disable Action = "disable"
	Start   Action = "start"
	Stop    Action = "stop"
	Restart Action = "restart"
	Status  Action = "status"
	Logs    Action = "logs"
)

// Actions lists every action in display order.
var Actions = []Action{Enable, Disable, Start, Stop, Restart, Status, Logs}

const (
	unitName        = "cloudflared"
	systemdMarker   = "/run/systemd/system"
	openRCBinary    = "rc-service"
	openRCUpdate    = "rc-update"
	openRCRunlevel  = "default"
	sysVBinary      = "service"
	sysVLogFile     = "/var/log/cloudflared.log"
	systemdNoPager  = "--no-pager"
	followLogsFlag  = "-f"
	journalUnitFlag = "-u"
)

// Manager runs service actions through the host init system.
type Manager struct {
	Name     string
	Runner   sysexec.Runner
	commands map[Action][]string
}

// Run executes action attached to the terminal.
func (m Manager) Run(ctx context.Context, action Action) error {
	args, ok := m.commands[action]
	if !ok {
		return fmt.Errorf("%s does not support %q", m.Name, action)
	}
	return m.Runner.Stream(ctx, sysexec.Command{Name: args[0], Args: args[1:], Privileged: true})
}

// Detect picks the manager for the running init system.
func Detect(r sysexec.Runner, exists func(string) bool, lookPath func(string) (string, error)) Manager {
	if exists(systemdMarker) {
		return systemd(r)
	}
	if _, err := lookPath(openRCBinary); err == nil {
		return openRC(r)
	}
	return sysV(r)
}

func systemd(r sysexec.Runner) Manager {
	return Manager{Name: "systemd", Runner: r, commands: map[Action][]string{
		Enable:  {"systemctl", "enable", unitName},
		Disable: {"systemctl", "disable", unitName},
		Start:   {"systemctl", "start", unitName},
		Stop:    {"systemctl", "stop", unitName},
		Restart: {"systemctl", "restart", unitName},
		Status:  {"systemctl", "status", systemdNoPager, unitName},
		Logs:    {"journalctl", journalUnitFlag, unitName, followLogsFlag},
	}}
}

func openRC(r sysexec.Runner) Manager {
	return Manager{Name: "openrc", Runner: r, commands: map[Action][]string{
		Enable:  {openRCUpdate, "add", unitName, openRCRunlevel},
		Disable: {openRCUpdate, "del", unitName, openRCRunlevel},
		Start:   {openRCBinary, unitName, "start"},
		Stop:    {openRCBinary, unitName, "stop"},
		Restart: {openRCBinary, unitName, "restart"},
		Status:  {openRCBinary, unitName, "status"},
		Logs:    {"tail", followLogsFlag, sysVLogFile},
	}}
}

func sysV(r sysexec.Runner) Manager {
	return Manager{Name: "sysvinit", Runner: r, commands: map[Action][]string{
		Start:   {sysVBinary, unitName, "start"},
		Stop:    {sysVBinary, unitName, "stop"},
		Restart: {sysVBinary, unitName, "restart"},
		Status:  {sysVBinary, unitName, "status"},
		Logs:    {"tail", followLogsFlag, sysVLogFile},
	}}
}

// unitFiles are where `cloudflared service install` registers the service.
var unitFiles = []string{
	"/etc/systemd/system/cloudflared.service",
	"/etc/init.d/cloudflared",
}

// Installed reports whether the service has been registered.
func Installed(exists func(string) bool) bool {
	for _, f := range unitFiles {
		if exists(f) {
			return true
		}
	}
	return false
}
