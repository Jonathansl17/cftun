package cli

import (
	"context"
	"io"

	"github.com/Jonathansl17/cftun/internal/cli/routecmd"
	"github.com/Jonathansl17/cftun/internal/cli/servicecmd"
	"github.com/Jonathansl17/cftun/internal/cli/tunnelcmd"
	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/dnsapi"
	"github.com/Jonathansl17/cftun/internal/health"
	"github.com/Jonathansl17/cftun/internal/hostfs"
	"github.com/Jonathansl17/cftun/internal/installer"
	"github.com/Jonathansl17/cftun/internal/quicktunnel"
	"github.com/Jonathansl17/cftun/internal/routes"
	"github.com/Jonathansl17/cftun/internal/service"
	"github.com/Jonathansl17/cftun/internal/store"
	"github.com/Jonathansl17/cftun/internal/sysexec"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type configLocator struct {
	Flag string
	Env  func(string) string
}

type infra struct {
	session     uikit.Session
	files       hostfs.FS
	runner      sysexec.Runner
	store       store.File
	tunnels     cloudflared.Client
	svc         service.Controller
	installer   installer.Installer
	tokens      dnsapi.TokenStore
	dns         dnsapi.Client
	health      health.Checker
	editor      routes.Editor
	validator   routes.Validator
	initializer routes.Initializer
	teardown    teardown.Procedure
	quick       quicktunnel.Runner
}

type groups struct {
	routes  *routecmd.Group
	tunnels *tunnelcmd.Group
	service *servicecmd.Group
}

type VersionSource interface {
	Version(ctx context.Context) (string, error)
}

type Gate struct {
	Session uikit.Session
	Tunnels VersionSource
	Install func(ctx context.Context) error
}

type App struct {
	uikit.Session
	Gate         Gate
	Files        hostfs.FS
	Store        store.File
	Tunnels      cloudflared.Client
	Installer    installer.Installer
	Tokens       dnsapi.TokenStore
	DNS          dnsapi.Client
	Teardown     teardown.Procedure
	Quick        quicktunnel.Runner
	Groups       groups
	CftunRemoved bool
}

type Reporter interface {
	Printf(format string, args ...any)
}

type teardownPrinter struct {
	Report Reporter
}

type entry struct {
	label    string
	path     []string
	children []entry
}

type menuScreen struct {
	title   string
	entries []entry
	leave   string
}
