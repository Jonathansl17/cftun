package service

import "github.com/Jonathansl17/cftun/internal/paths"

const (
	Enable  Action = "enable"
	Disable Action = "disable"
	Start   Action = "start"
	Stop    Action = "stop"
	Restart Action = "restart"
	Status  Action = "status"
	Logs    Action = "logs"
)

const (
	unitName                = paths.CloudflaredBinary
	systemdName             = "systemd"
	openRCName              = "openrc"
	sysVName                = "sysvinit"
	systemctlBinary         = "systemctl"
	journalctlBinary        = "journalctl"
	openRCBinary            = "rc-service"
	openRCUpdateBinary      = "rc-update"
	openRCRunlevel          = "default"
	openRCAddVerb           = "add"
	openRCDeleteVerb        = "del"
	sysVBinary              = "service"
	tailBinary              = "tail"
	systemdNoPagerFlag      = "--no-pager"
	followLogsFlag          = "-f"
	journalUnitFlag         = "-u"
	unsupportedActionFormat = "%s does not support %q"
)

var Actions = []Action{Enable, Disable, Start, Stop, Restart, Status, Logs}

var unitFiles = []string{paths.ServiceUnit, paths.InitScript}
