package authcmd

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
)

type Cloudflared interface {
	Login(ctx context.Context) error
}

type Files interface {
	Exists(path string) bool
}

type Home interface {
	Home() (string, error)
}

type Tokens interface {
	Load() (string, error)
	Save(token string) error
	Clear() error
	Location() string
}

type Deps struct {
	Session     uikit.Session
	Cloudflared Cloudflared
	Files       Files
	Home        Home
	Tokens      Tokens
	Env         func(string) string
}

type Group struct {
	deps Deps
}
