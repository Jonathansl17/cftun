package cli

import (
	"time"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

const (
	binaryName    = "cftun"
	flagConfig    = "config"
	configPathEnv = "CFTUN_CONFIG"
	noDefault     = ""

	exitOK      = 0
	exitFailure = 1

	apiTimeout      = 15 * time.Second
	probeTimeout    = 5 * time.Second
	downloadTimeout = 5 * time.Minute

	resolvePathFormat = "resolve path: %w"
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

var mainMenu = []entry{
	{label: msg.MenuSetup, path: []string{"setup"}},
	{label: msg.MenuQuick, path: []string{"tunnel"}},
	{label: msg.MenuRoutes, children: []entry{
		{label: msg.MenuList, path: []string{"list"}},
		{label: msg.MenuAdd, path: []string{"add"}},
		{label: msg.MenuEdit, path: []string{"edit"}},
		{label: msg.MenuRemove, path: []string{"rm"}},
		{label: msg.MenuCheck, path: []string{"check"}},
		{label: msg.MenuValidate, path: []string{"validate"}},
	}},
	{label: msg.MenuService, children: []entry{
		{label: msg.MenuServiceInstall, path: []string{"service", "install"}},
		{label: msg.MenuServiceStatus, path: []string{"service", "status"}},
		{label: msg.MenuServiceRestart, path: []string{"service", "restart"}},
		{label: msg.MenuServiceStart, path: []string{"service", "start"}},
		{label: msg.MenuServiceStop, path: []string{"service", "stop"}},
		{label: msg.MenuServiceEnable, path: []string{"service", "enable"}},
		{label: msg.MenuServiceDisable, path: []string{"service", "disable"}},
		{label: msg.MenuServiceLogs, path: []string{"service", "logs"}},
	}},
	{label: msg.MenuTunnel, children: []entry{
		{label: msg.MenuTunnelList, path: []string{"tunnels", "list"}},
		{label: msg.MenuTunnelCreate, path: []string{"tunnels", "create"}},
		{label: msg.MenuTunnelDelete, path: []string{"tunnels", "delete"}},
		{label: msg.MenuInit, path: []string{"init"}},
	}},
	{label: msg.MenuAccount, children: []entry{
		{label: msg.MenuInstall, path: []string{"install"}},
		{label: msg.MenuLogin, path: []string{"login"}},
		{label: msg.MenuTokenSet, path: []string{"token", "set"}},
		{label: msg.MenuTokenStatus, path: []string{"token", "status"}},
		{label: msg.MenuTokenClear, path: []string{"token", "clear"}},
	}},
	{label: msg.MenuUninstall, path: []string{"uninstall"}},
}
