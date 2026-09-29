package cli

import (
	"time"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
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

	noPort = ""
)

var menuGroups = []menuGroupInfo{
	{group: uikit.GroupRoutes, label: msg.MenuRoutes, order: uikit.OrderRoutes},
	{group: uikit.GroupService, label: msg.MenuService, order: uikit.OrderService},
	{group: uikit.GroupTunnels, label: msg.MenuTunnel, order: uikit.OrderTunnels},
	{group: uikit.GroupAccount, label: msg.MenuAccount, order: uikit.OrderAccount},
}
