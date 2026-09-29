package msg

// Interactive menu labels.
const (
	MenuTitle = "\ncftun - what do you want to do?"
	MenuExit  = "Exit"
	MenuBack  = "Back"

	MenuSetup     = "Guided setup (start here)"
	MenuQuick     = "Temporary tunnel (no domain needed)"
	MenuRoutes    = "Routes (hostname -> port)"
	MenuService   = "Service"
	MenuTunnel    = "Named tunnels"
	MenuAccount   = "cloudflared, login and API token"
	MenuUninstall = "Uninstall everything"

	MenuList     = "List routes"
	MenuAdd      = "Add route"
	MenuEdit     = "Edit route"
	MenuRemove   = "Remove route"
	MenuCheck    = "Check routes (local and public)"
	MenuValidate = "Validate config"

	MenuServiceInstall = "Install and start service"
	MenuServiceStatus  = "Status"
	MenuServiceRestart = "Restart"
	MenuServiceStart   = "Start"
	MenuServiceStop    = "Stop"
	MenuServiceEnable  = "Enable on boot"
	MenuServiceDisable = "Disable on boot"
	MenuServiceLogs    = "Follow logs"

	MenuTunnelList   = "List tunnels"
	MenuTunnelCreate = "Create tunnel"
	MenuTunnelDelete = "Delete tunnel"
	MenuInit         = "Write config for a tunnel"

	MenuInstall     = "Install cloudflared"
	MenuLogin       = "Login to Cloudflare"
	MenuTokenSet    = "Set API token"
	MenuTokenStatus = "API token status"
	MenuTokenClear  = "Clear API token"
)

// Setup wizard steps.
const (
	SetupStepFormat  = "\n==> [%d/%d] %s\n"
	SetupStepInstall = "Install cloudflared"
	SetupStepLogin   = "Login to Cloudflare"
	SetupStepTunnel  = "Tunnel and config file"
	SetupStepService = "System service"
	SetupStepToken   = "API token (optional)"
	SetupStepRoutes  = "Routes"
)
