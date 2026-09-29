package service

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func (c Controller) Run(ctx context.Context, action Action) error {
	args, ok := c.commands[action]
	if !ok {
		return &UnsupportedActionError{Name: c.Name, Action: action}
	}
	return c.Runner.Stream(ctx, sysexec.Command{Name: args[0], Args: args[1:], Privileged: true})
}

func (c Controller) Restart(ctx context.Context) error {
	return c.Run(ctx, Restart)
}

func (c Controller) Stop(ctx context.Context) error {
	return c.Run(ctx, Stop)
}

func (c Controller) Registered() bool {
	for _, file := range unitFiles {
		if c.exists(file) {
			return true
		}
	}
	return false
}
