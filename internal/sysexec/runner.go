package sysexec

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func (c Command) String() string {
	return strings.Join(append([]string{c.Name}, c.Args...), commandStringSpacing)
}

func NewShell() *Shell {
	return &Shell{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		IsRoot: func() bool { return os.Geteuid() == 0 },
	}
}

func (s *Shell) Output(ctx context.Context, cmd Command) (string, error) {
	var stdout, stderr bytes.Buffer
	c := s.build(ctx, cmd)
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		return "", &CommandError{Command: cmd, Stderr: strings.TrimSpace(stderr.String()), Err: err}
	}
	return stdout.String(), nil
}

func (s *Shell) Stream(ctx context.Context, cmd Command) error {
	c := s.build(ctx, cmd)
	c.Stdin, c.Stdout, c.Stderr = s.Stdin, s.Stdout, s.Stderr
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	defer signal.Stop(interrupts)
	if err := c.Run(); err != nil {
		return &CommandError{Command: cmd, Err: err}
	}
	return nil
}

func (s *Shell) build(ctx context.Context, cmd Command) *exec.Cmd {
	resolved := s.resolve(cmd)
	return exec.CommandContext(ctx, resolved.Name, resolved.Args...)
}

func (s *Shell) resolve(cmd Command) Command {
	if !cmd.Privileged || s.IsRoot() {
		return cmd
	}
	return Command{Name: elevator, Args: append([]string{cmd.Name}, cmd.Args...)}
}
