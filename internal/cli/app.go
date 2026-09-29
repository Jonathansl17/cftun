package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/dnsapi"
	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/hostfs"
	"github.com/Jonathansl17/cftun/internal/installer"
	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/prompt"
	"github.com/Jonathansl17/cftun/internal/quicktunnel"
	"github.com/Jonathansl17/cftun/internal/routes"
	"github.com/Jonathansl17/cftun/internal/service"
	"github.com/Jonathansl17/cftun/internal/store"
	"github.com/Jonathansl17/cftun/internal/sysexec"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

const (
	ConfigPathEnv   = "CFTUN_CONFIG"
	apiTimeout      = 15 * time.Second
	probeTimeout    = 5 * time.Second
	downloadTimeout = 5 * time.Minute
)

type App struct {
	Out         io.Writer
	Home        string
	Files       hostfs.FS
	ConfigPath  string
	Runner      sysexec.Runner
	Store       store.File
	Tunnels     cloudflared.Client
	Service     service.Controller
	Installer   installer.Installer
	Tokens      dnsapi.TokenStore
	DNS         dnsapi.Client
	Health      health.Checker
	Prompt      prompt.Prompter
	Editor      routes.Editor
	Checker     routes.Validator
	Initializer routes.Initializer
	Teardown    teardown.Procedure
	Quick       quicktunnel.Runner

	CftunRemoved bool
}

func NewApp(configPath string, in io.Reader, out io.Writer) (*App, error) {
	files := hostfs.FS{}
	home, err := files.Home()
	if err != nil {
		return nil, err
	}
	tokens, err := dnsapi.NewTokenStore()
	if err != nil {
		return nil, err
	}
	a := &App{Out: out, Home: home, Files: files, ConfigPath: configPath, Runner: sysexec.NewShell(), Tokens: tokens}
	a.wire(in)
	return a, nil
}

func (a *App) wire(in io.Reader) {
	a.Store = store.File{Locator: FixedLocator{Path: a.ConfigPath}, Writer: store.ElevatedWriter{Runner: a.Runner}}
	a.Tunnels = cloudflared.Client{Runner: a.Runner}
	a.Quick = quicktunnel.Runner{Tunnels: a.Tunnels, Files: a.Files}
	a.Service = service.Detect(a.Runner, a.Files.Exists, exec.LookPath)
	a.Installer = installer.Installer{
		Runner:     a.Runner,
		Downloader: installer.HTTPDownloader{Client: &http.Client{Timeout: downloadTimeout}},
		Releases:   installer.GitHubReleases{Client: &http.Client{Timeout: apiTimeout}},
		Families:   installer.OSReleaseSource{},
		GoArch:     runtime.GOARCH,
	}
	a.DNS = dnsapi.Client{HTTP: &http.Client{Timeout: apiTimeout}, BaseURL: dnsapi.DefaultBaseURL, Tokens: a.Tokens}
	a.Health = health.Checker{HTTP: &http.Client{Timeout: probeTimeout}}
	a.Prompt = prompt.NewConsole(in, a.Out, promptTexts())
	a.wireRoutes()
}

func (a *App) wireRoutes() {
	dns := routes.DNSCleaner{API: a.DNS}
	a.Checker = routes.Validator{Path: a.Store, Tunnels: a.Tunnels}
	a.Editor = routes.Editor{Store: a.Store, Checker: a.Checker, Router: a.Tunnels, DNS: dns, Restarter: a.Service}
	a.Initializer = routes.Initializer{Store: a.Store, Checker: a.Checker, Catalog: a.Tunnels, Home: a.Files, Files: a.Files}
	a.Teardown = teardown.Procedure{
		Config: a.Store, Tunnels: a.Tunnels, Service: a.Service, Package: a.Installer, DNS: dns,
		Files:    teardown.SafeRemover{Runner: a.Runner, Inspector: a.Files},
		Temp:     a.Files,
		Observer: teardownPrinter{Report: a},
	}
}

func (a *App) Printf(format string, args ...any) {
	fmt.Fprintf(a.Out, format, args...)
}

func (a *App) UserCloudflaredDir() string {
	return filepath.Join(a.Home, paths.UserDirName)
}

func ResolveConfigPath(flag string) string {
	if flag != "" {
		return flag
	}
	if env := os.Getenv(ConfigPathEnv); env != "" {
		return env
	}
	return paths.DefaultConfig
}
