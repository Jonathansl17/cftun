package sysexec

import (
	"context"
	"io"
)

type Command struct {
	Name       string
	Args       []string
	Privileged bool
}

type Runner interface {
	Output(ctx context.Context, cmd Command) (string, error)
	Stream(ctx context.Context, cmd Command) error
}

type Shell struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	IsRoot func() bool
}

type CommandError struct {
	Command Command
	Stderr  string
	Err     error
}
