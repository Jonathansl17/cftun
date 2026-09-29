package msg

import "fmt"

const (
	PromptFormat              = "%s: "
	OptionFormat              = "  %d) %s\n"
	ConfirmFormat             = "%s [y/N]: "
	PromptChoice              = "Choose an option"
	PromptHostname            = "Public hostname (e.g. api.example.com)"
	PromptPort                = "Local port"
	PromptNewHostname         = "New hostname"
	PromptNewPort             = "New port"
	PromptPickRule            = "Which route?"
	PromptPickTunnel          = "Which tunnel?"
	PromptTunnelName          = "Tunnel name"
	PromptToken               = "Cloudflare API token"
	RuleArrow                 = " -> "
	ConfirmUninstall          = "This deletes the tunnel, its DNS records, the service, the package, all credentials and cftun itself. Continue?"
	ConfirmDeleteTunnel       = "Delete tunnel %q?"
	ConfirmOverwrite          = "A config file already exists. Overwrite it?"
	ConfirmSetToken           = "Save an API token now?"
	ConfirmAddRoute           = "Add a route?"
	ConfirmHasDomain          = "Do you have a domain added to your Cloudflare account?"
	ConfirmInstallCloudflared = "cloudflared is not installed. Install it now?"
)

func CurrentValue(current string) string {
	return fmt.Sprintf(" [%s]", current)
}

const (
	InfoAdded            = "Route ready: https://%s -> %s\n"
	InfoRemoved          = "Route removed: %s\n"
	InfoNoRules          = "No routes yet. Add one with \"cftun add\".\n"
	InfoNoTunnels        = "No tunnels yet. Create one with \"cftun tunnels create\".\n"
	InfoDNSDeleted       = "Deleted %d DNS record(s) for %s\n"
	InfoValidated        = "%s"
	InfoRestartSkipped   = "Service not restarted; run \"cftun service restart\" to apply.\n"
	InfoAlreadyInstalled = "cloudflared already installed: %s\n"
	InfoInstalling       = "Installing cloudflared for the %s family...\n"
	InfoInstalled        = "Installed: %s\n"
	InfoAlreadyLoggedIn  = "Already logged in (%s)\n"
	InfoTunnelReused     = "Tunnel %q already exists, reusing it\n"
	InfoConfigWritten    = "Config written to %s\n"
	InfoConfigKept       = "Keeping existing config %s\n"
	InfoServiceExists    = "Service already registered\n"
	InfoTokenSet         = "API token configured\n"
	InfoTokenMissing     = "No API token (set one with \"cftun token set\" or the %s variable)\n"
	InfoTokenSaved       = "Token saved to %s\n"
	InfoSetupDone        = "\nDone. Manage routes any time with \"cftun\".\n"
	InfoRemoving         = "Removing %s\n"
	InfoNoDomain         = "Without a domain, cftun starts a temporary tunnel. Add a domain to Cloudflare later and run \"cftun setup\" again for a permanent one.\n"
	InfoQuickStarting    = "Starting a temporary tunnel to %s. Wait for the https://<name>.trycloudflare.com line below; press your interrupt key (Ctrl+C) to stop.\n\n"
	InfoQuickStopped     = "\nTemporary tunnel stopped.\n"
	InfoUninstalled      = "cloudflared and all its files were removed.\n"
	InfoCftunRemoved     = "cftun was removed.\n"
	WarnManualDNS        = "No API token: delete the CNAME record of %s in the Cloudflare dashboard (DNS > Records).\n"
	WarnNoConfig         = "Skipping routes and tunnel: %v\n"
	WarnStepFailed       = "warning: %s failed: %v\n"
)

const (
	ListHeader       = "HOSTNAME\tSERVICE\tLOCAL"
	TunnelHeader     = "NAME\tID"
	CheckTitle       = "\n%s\n"
	CheckLine        = "  %-40s %s\n"
	CheckHint        = "    hint: %s\n"
	StatusDown       = "DOWN"
	StatusOKFormat   = "OK (%d)"
	StatusFailFormat = "FAIL (%d)"
)
