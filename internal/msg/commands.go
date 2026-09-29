package msg

const (
	RootShort = "Manage a Cloudflare Tunnel and its public hostnames"
	RootLong  = `cftun installs cloudflared, creates a tunnel, publishes local ports on
your domain and keeps the service running, without editing files by hand.

Run it without arguments for an interactive menu, or use the subcommands
with flags for scripts. Start with "cftun setup" on a new machine.`

	SetupShort = "Guided setup: install, login, tunnel, config, service and routes"
	SetupLong  = `Runs every step needed to go from nothing to a working tunnel. Steps
already done are detected and skipped, so it is safe to run it again.
Without a domain on Cloudflare, it starts a temporary tunnel instead.`

	InstallShort   = "Install cloudflared using the distribution's native method"
	UninstallShort = "Remove cloudflared and everything cftun created"
	UninstallLong  = `Stops and unregisters the service, deletes the DNS records of every route
(needs an API token), deletes the tunnel, removes the package, the config,
its backup, the credentials in ~/.cloudflared and the saved API token.`

	LoginShort       = "Authorize cloudflared with your Cloudflare account"
	TokenShort       = "Manage the API token used to delete DNS records"
	TokenSetShort    = "Save an API token"
	TokenClearShort  = "Forget the saved API token"
	TokenStatusShort = "Show whether an API token is configured"
	TokenLong        = `cloudflared can create DNS records but not delete them. To let cftun
remove the CNAME of a route, create a token at
https://dash.cloudflare.com/profile/api-tokens with the "Edit zone DNS"
template and save it with "cftun token set". Without it, cftun tells you
which records to delete in the dashboard.
`

	TunnelShort       = "Create, list and delete named tunnels (needs a domain)"
	TunnelCreateShort = "Create a tunnel and its credentials file"
	TunnelListShort   = "List the tunnels of the account"
	TunnelDeleteShort = "Delete a tunnel"
	InitShort         = "Write the config file for an existing tunnel"

	QuickShort = "Temporary public URL for a local port, no domain or account needed"
	QuickLong  = `Starts a Cloudflare Quick Tunnel: a random https://<name>.trycloudflare.com
URL that forwards to a local port. It needs no domain, account or login.

The URL changes every run and lives only while the command runs. One port
per run, not a system service, no uptime guarantee, about 200 concurrent
requests and no Server-Sent Events. Use a named tunnel ("cftun setup") for
anything permanent.`

	AddShort      = "Publish a local port on a hostname"
	RemoveShort   = "Unpublish a hostname"
	EditShort     = "Change the hostname or port of a route"
	ListShort     = "List routes and whether their local port answers"
	CheckShort    = "Test every route locally and from the internet"
	ValidateShort = "Validate the ingress rules with cloudflared"

	ServiceShort        = "Control the cloudflared system service"
	ServiceInstallShort = "Register, enable and start the service"
)

var ServiceActionShort = map[string]string{
	"enable":  "Start the service on boot",
	"disable": "Do not start the service on boot",
	"start":   "Start the service",
	"stop":    "Stop the service",
	"restart": "Restart the service to apply config changes",
	"status":  "Show the service status",
	"logs":    "Follow the service logs (Ctrl+C to stop)",
}

const (
	FlagConfig       = "cloudflared config file (default /etc/cloudflared/config.yml, env CFTUN_CONFIG)"
	FlagHost         = "public hostname, e.g. api.example.com"
	FlagNewHost      = "new public hostname"
	FlagPort         = "local port the app listens on"
	FlagNoDNS        = "do not create or delete the DNS record"
	FlagNoRestart    = "do not restart the service after the change"
	FlagTunnel       = "tunnel name or UUID"
	FlagForceInit    = "overwrite an existing config without asking"
	FlagForceInstall = "reinstall even if cloudflared is present"
	FlagYes          = "do not ask for confirmation"
)
