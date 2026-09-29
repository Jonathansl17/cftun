package cli

import (
	"net/http"
	"os/exec"
	"runtime"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/dnsapi"
	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/hostfs"
	"github.com/Jonathansl17/cftun/internal/installer"
	"github.com/Jonathansl17/cftun/internal/prompt"
	"github.com/Jonathansl17/cftun/internal/quicktunnel"
	"github.com/Jonathansl17/cftun/internal/routes"
	"github.com/Jonathansl17/cftun/internal/service"
	"github.com/Jonathansl17/cftun/internal/store"
	"github.com/Jonathansl17/cftun/internal/sysexec"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

func newInfra(streams Streams, locator *configLocator) (infra, error) {
	tokens, err := dnsapi.NewTokenStore()
	if err != nil {
		return infra{}, err
	}
	files := hostfs.FS{}
	runner := sysexec.NewShell()
	console := prompt.NewConsole(streams.In, streams.Out, promptTexts())
	in := infra{
		session: uikit.Session{Out: streams.Out, Prompt: console},
		files:   files,
		runner:  runner,
		store:   store.File{Locator: locator, Writer: store.ElevatedWriter{Runner: runner}},
		tunnels: cloudflared.Client{Runner: runner},
		svc:     service.Detect(runner, files.Exists, exec.LookPath),
		tokens:  tokens,
		health:  health.Checker{HTTP: &http.Client{Timeout: probeTimeout}},
	}
	in.installer = newInstaller(runner)
	in.dns = dnsapi.Client{HTTP: &http.Client{Timeout: apiTimeout}, BaseURL: dnsapi.DefaultBaseURL, Tokens: tokens}
	in.quick = quicktunnel.Runner{Tunnels: in.tunnels, Files: files}
	in.wireRoutes()
	return in, nil
}

func newInstaller(runner sysexec.Runner) installer.Installer {
	return installer.Installer{
		Runner:     runner,
		Downloader: installer.HTTPDownloader{Client: &http.Client{Timeout: downloadTimeout}},
		Releases:   installer.GitHubReleases{Client: &http.Client{Timeout: apiTimeout}},
		Families:   installer.OSReleaseSource{},
		GoArch:     runtime.GOARCH,
	}
}

func (in *infra) wireRoutes() {
	dns := routes.DNSCleaner{API: in.dns}
	in.validator = routes.Validator{Path: in.store, Tunnels: in.tunnels}
	in.editor = routes.Editor{Store: in.store, Checker: in.validator, Router: in.tunnels, DNS: dns, Restarter: in.svc}
	in.initializer = routes.Initializer{
		Store: in.store, Checker: in.validator, Catalog: in.tunnels, Home: in.files, Files: in.files,
	}
	in.teardown = teardown.Procedure{
		Config: in.store, Tunnels: in.tunnels, Service: in.svc, Package: in.installer, DNS: dns,
		Files:    teardown.SafeRemover{Runner: in.runner, Inspector: in.files},
		Temp:     in.files,
		Observer: teardownPrinter{Report: in.session},
	}
}
