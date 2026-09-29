package paths

const (
	CloudflaredBinary = "cloudflared"
	UserDirName       = ".cloudflared"
	CertFile          = "cert.pem"

	ConfigDir      = "/etc/cloudflared"
	DefaultConfig  = "/etc/cloudflared/config.yml"
	LocalConfigDir = "/usr/local/etc/cloudflared"
	RootHome       = "/root"
	RootUserDir    = "/root/.cloudflared"

	LogFile      = "/var/log/cloudflared.log"
	ErrorLogFile = "/var/log/cloudflared.err"

	ServiceUnit       = "/etc/systemd/system/cloudflared.service"
	UpdateServiceUnit = "/etc/systemd/system/cloudflared-update.service"
	UpdateTimerUnit   = "/etc/systemd/system/cloudflared-update.timer"
	InitScript        = "/etc/init.d/cloudflared"

	OSRelease     = "/etc/os-release"
	SystemdMarker = "/run/systemd/system"
	InstallPath   = "/usr/local/bin/cloudflared"

	BashCompletion = "/usr/share/bash-completion/completions/cftun"
	ZshCompletion  = "/usr/share/zsh/site-functions/_cftun"
	FishCompletion = "/usr/share/fish/vendor_completions.d/cftun.fish"
)
