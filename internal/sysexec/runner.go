// Package sysexec runs external programs, elevating with sudo when required.
package sysexec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// elevator is the program used to gain root privileges.
const (
	elevator            = "sudo"
	exitCodeInterrupted = 130
)

// Command describes one external program invocation.
type Command struct {
	Name       string
	Args       []string
	Privileged bool
}

// String renders the command the way a user would type it.
func (c Command) String() string {
	return strings.Join(append([]string{c.Name}, c.Args...), " ")
}

// Runner executes commands.
type Runner interface {
	// Output runs the command and returns its captured stdout.
	Output(ctx context.Context, cmd Command) (string, error)
	// Stream runs the command attached to the terminal.
	Stream(ctx context.Context, cmd Command) error
}

// Shell is the Runner backed by the operating system.
type Shell struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	IsRoot func() bool
}

// NewShell returns a Shell wired to the process standard streams.
func NewShell() *Shell {
	return &Shell{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		IsRoot: func() bool { return os.Geteuid() == 0 },
	}
}

// Output implements Runner.
func (s *Shell) Output(ctx context.Context, cmd Command) (string, error) {
	var stdout, stderr bytes.Buffer
	c := s.build(ctx, cmd)
	c.Stdout = &stdout
	c.Stderr = &stderr
	if err := c.Run(); err != nil {
		return "", fmt.Errorf("%s: %w: %s", cmd, err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

// Stream implements Runner.
func (s *Shell) Stream(ctx context.Context, cmd Command) error {
	c := s.build(ctx, cmd)
	c.Stdin, c.Stdout, c.Stderr = s.Stdin, s.Stdout, s.Stderr
	// Ctrl+C belongs to the child (e.g. following logs); keep cftun alive.
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	defer signal.Stop(interrupts)
	if err := c.Run(); err != nil {
		return fmt.Errorf("%s: %w", cmd, err)
	}
	return nil
}

func (s *Shell) build(ctx context.Context, cmd Command) *exec.Cmd {
	resolved := s.resolve(cmd)
	return exec.CommandContext(ctx, resolved.Name, resolved.Args...)
}

// resolve prefixes privileged commands with the elevator when not root.
func (s *Shell) resolve(cmd Command) Command {
	if !cmd.Privileged || s.IsRoot() {
		return cmd
	}
	return Command{Name: elevator, Args: append([]string{cmd.Name}, cmd.Args...)}
}

// Interrupted reports whether err is a child exiting because of Ctrl+C.
func Interrupted(err error) bool {
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		return false
	}
	status, ok := exit.Sys().(syscall.WaitStatus)
	return ok && (status.Signaled() && status.Signal() == syscall.SIGINT || status.ExitStatus() == exitCodeInterrupted)
}
