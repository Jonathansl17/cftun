package msg

const (
	ErrRolledBack         = "cloudflared rejected the config, previous version restored"
	ErrDNSRolledBack      = "could not create the DNS route, previous config restored"
	ErrTunnelNotFound     = "tunnel not found, create it with `cftun tunnels create`"
	ErrCredentialsMissing = "credentials file not found"
	ErrUnsafeTarget       = "refusing to remove %s: only regular config files and cloudflared directories can be removed"
)
