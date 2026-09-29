package msg

const (
	ErrorFormat           = "error: %v\n"
	ErrChoiceRange        = "enter a number between 1 and %d"
	ErrEmpty              = "a value is required"
	ErrNoRules            = "there are no routes yet, add one first"
	ErrCancelled          = "cancelled"
	ErrCloudflaredMissing = "cloudflared is required, install it with \"cftun install\""
)

const (
	ErrRolledBack         = "cloudflared rejected the config, previous version restored"
	ErrDNSRolledBack      = "could not create the DNS route, previous config restored"
	ErrTunnelNotFound     = "tunnel not found, create it with `cftun tunnels create`"
	ErrCredentialsMissing = "credentials file not found"
	ErrUnsafeTarget       = "refusing to remove %s: only regular config files and cloudflared directories can be removed"
	ErrMalformed          = "config must be a mapping with an ingress list"
	ErrDuplicate          = "hostname already has a rule"
	ErrNotFound           = "hostname has no rule"
	ErrInvalidHostname    = "invalid hostname, expected something like app.example.com"
	ErrInvalidPortFormat  = "invalid port, expected a number between %d and %d"
	ErrConfigMissing      = "config file not found, run `cftun init` first"
)
