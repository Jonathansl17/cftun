package sysexec

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

func (e *CommandError) Error() string {
	text := fmt.Sprintf(commandErrorFormat, e.Command, e.Err)
	if e.Stderr == "" {
		return text
	}
	return fmt.Sprintf(commandStderrFormat, text, e.Stderr)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

func ExitCode(err error) (int, bool) {
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		return 0, false
	}
	return exit.ExitCode(), true
}

func Interrupted(err error) bool {
	return signaledByInterrupt(err) || exitedInterrupted(err)
}

func signaledByInterrupt(err error) bool {
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		return false
	}
	status, ok := exit.Sys().(syscall.WaitStatus)
	return ok && status.Signaled() && status.Signal() == syscall.SIGINT
}

func exitedInterrupted(err error) bool {
	code, ok := ExitCode(err)
	return ok && code == interruptedExitCode
}
