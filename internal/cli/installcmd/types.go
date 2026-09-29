package installcmd

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/installer"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

type Installer interface {
	Install(ctx context.Context) error
	Family() (installer.Family, error)
}

type Cloudflared interface {
	Version(ctx context.Context) (string, error)
}

type Procedure interface {
	Run(ctx context.Context, params teardown.Params) error
}

type Tokens interface {
	Clear() error
}

type Config interface {
	Path() string
	BackupPath() string
}

type Home interface {
	Home() (string, error)
}

type Files interface {
	Executable() (string, error)
	Exists(path string) bool
}

type Remover interface {
	RemoveFile(ctx context.Context, path string) error
}

type Deps struct {
	Session     uikit.Session
	Installer   Installer
	Cloudflared Cloudflared
	Procedure   Procedure
	Tokens      Tokens
	Config      Config
	Home        Home
	Files       Files
	Remover     Remover
}

type Group struct {
	deps    Deps
	removed bool
}

type Observer struct {
	Session uikit.Session
}
