package cli

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/authcmd"
	"github.com/Jonathansl17/cftun/internal/cli/installcmd"
	"github.com/Jonathansl17/cftun/internal/cli/quickcmd"
	"github.com/Jonathansl17/cftun/internal/cli/routecmd"
	"github.com/Jonathansl17/cftun/internal/cli/servicecmd"
	"github.com/Jonathansl17/cftun/internal/cli/setupcmd"
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
	setup   *setupcmd.Group
	install *installcmd.Group
	auth    *authcmd.Group
	quick   *quickcmd.Group
	tunnels *tunnelcmd.Group
	routes  *routecmd.Group
	service *servicecmd.Group
}

type commandProvider interface {
	Commands() []*cobra.Command
}

type VersionSource interface {
	Version(ctx context.Context) (string, error)
}

type Gate struct {
	Session uikit.Session
	Tunnels VersionSource
	Install func(ctx context.Context) error
}

type menu struct {
	session  uikit.Session
	gate     Gate
	finished func() bool
}

type menuNode struct {
	label    string
	order    int
	cmd      *cobra.Command
	children []menuNode
}

type menuScreen struct {
	title string
	nodes []menuNode
	leave string
}

type menuGroupInfo struct {
	group uikit.MenuGroup
	label string
	order int
}
