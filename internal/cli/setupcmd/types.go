package setupcmd

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
)

type Step struct {
	Title string
	Run   func(ctx context.Context) error
}

type Deps struct {
	Session         uikit.Session
	Steps           []Step
	TemporaryTunnel func(ctx context.Context) error
}

type Group struct {
	deps Deps
}
