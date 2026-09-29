// Package cli defines the cftun commands.
package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Jonathansl17/cloudfare/internal/cloudflared"
	"github.com/Jonathansl17/cloudfare/internal/dnsapi"
	"github.com/Jonathansl17/cloudfare/internal/health"
	"github.com/Jonathansl17/cloudfare/internal/installer"
	"github.com/Jonathansl17/cloudfare/internal/prompt"
	"github.com/Jonathansl17/cloudfare/internal/routes"
	"github.com/Jonathansl17/cloudfare/internal/service"
	"github.com/Jonathansl17/cloudfare/internal/store"
	"github.com/Jonathansl17/cloudfare/internal/sysexec"
)

const (
	// DefaultConfigPath is where the cloudflared service reads its config.
	DefaultConfigPath = "/etc/cloudflared/config.yml"
	// ConfigPathEnv overrides the config path.
	ConfigPathEnv   = "CFTUN_CONFIG"
	apiTimeout      = 15 * time.Second
	probeTimeout    = 5 * time.Second
	downloadTimeout = 5 * time.Minute
)

// App holds every collaborator the commands use.
type App struct {
	Out        io.Writer
	Home       string
	ConfigPath string
	Runner     sysexec.Runner
	Store      store.File
	Tunnels    cloudflared.Client
	Service    service.Manager
	Installer  installer.Installer
	Tokens     dnsapi.TokenStore
	DNS        dnsapi.Client
	Health     health.Checker
	Prompt     prompt.Prompter
	Routes     routes.Manager
}

// NewApp wires the real implementations for configPath.
func NewApp(configPath string, in io.Reader, out io.Writer) (*App, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("locate home: %w", err)
	}
	tokens, err := dnsapi.NewTokenStore()
	if err != nil {
		return nil, err
	}
	token, err := tokens.Load()
	if err != nil {
		return nil, err
	}
	family, err := detectFamily()
	if err != nil {
		return nil, err
	}
	a := &App{Out: out, Home: home, ConfigPath: configPath, Runner: sysexec.NewShell(), Tokens: tokens}
	a.wire(family, token, in)
	return a, nil
}

func (a *App) wire(family installer.Family, token string, in io.Reader) {
	a.Store = store.File{Path: a.ConfigPath, Writer: store.ElevatedWriter{Runner: a.Runner}}
	a.Tunnels = cloudflared.Client{Runner: a.Runner}
	a.Service = service.Detect(a.Runner, fileExists, exec.LookPath)
	a.Installer = installer.Installer{
		Runner:     a.Runner,
		Downloader: installer.HTTPDownloader{Client: &http.Client{Timeout: downloadTimeout}},
		Family:     family,
		GoArch:     runtime.GOARCH,
	}
	a.DNS = dnsapi.Client{HTTP: &http.Client{Timeout: apiTimeout}, BaseURL: dnsapi.DefaultBaseURL, Token: token}
	a.Health = health.Checker{HTTP: &http.Client{Timeout: probeTimeout}}
	a.Prompt = prompt.NewConsole(in, a.Out)
	a.Routes = routes.Manager{
		ConfigPath: a.ConfigPath, Store: a.Store, Tunnels: a.Tunnels,
		Service: a.Service, DNS: a.DNS, Report: a,
	}
}

// Printf implements routes.Reporter.
func (a *App) Printf(format string, args ...any) {
	fmt.Fprintf(a.Out, format, args...)
}

// UserCloudflaredDir is the per-user cloudflared directory.
func (a *App) UserCloudflaredDir() string {
	return filepath.Join(a.Home, cloudflared.HomeDir)
}

func detectFamily() (installer.Family, error) {
	f, err := os.Open(installer.OSReleasePath)
	if errors.Is(err, fs.ErrNotExist) {
		return installer.FamilyUnknown, nil
	}
	if err != nil {
		return installer.FamilyUnknown, fmt.Errorf("read %s: %w", installer.OSReleasePath, err)
	}
	defer f.Close()
	return installer.DetectFamily(f)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ResolveConfigPath applies flag, then environment, then default.
func ResolveConfigPath(flag string) string {
	if flag != "" {
		return flag
	}
	if env := os.Getenv(ConfigPathEnv); env != "" {
		return env
	}
	return DefaultConfigPath
}
