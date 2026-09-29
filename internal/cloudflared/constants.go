package cloudflared

const (
	credentialsExt = ".json"

	subTunnel    = "tunnel"
	subLogin     = "login"
	subCreate    = "create"
	subCleanup   = "cleanup"
	subDelete    = "delete"
	subList      = "list"
	subRoute     = "route"
	subDNS       = "dns"
	subIngress   = "ingress"
	subValidate  = "validate"
	subService   = "service"
	subInstall   = "install"
	subUninstall = "uninstall"

	flagVersion      = "--version"
	flagForce        = "-f"
	flagOutput       = "--output"
	flagConfig       = "--config"
	flagNoAutoupdate = "--no-autoupdate"
	flagGracePeriod  = "--grace-period"
	flagURL          = "--url"

	outputJSON       = "json"
	quickGracePeriod = "1s"

	parseTunnelListFormat = "parse tunnel list: %w"
)
