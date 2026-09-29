package quickcmd

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
)

type Ensurer interface {
	EnsureInstalled(ctx context.Context, force bool) error
}

type Tunnel interface {
	Run(ctx context.Context, localURL string) error
}

type Deps struct {
	Session uikit.Session
	Ensurer Ensurer
	Tunnel  Tunnel
}

type Group struct {
	deps Deps
}
