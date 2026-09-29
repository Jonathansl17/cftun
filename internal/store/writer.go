package store

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func (w ElevatedWriter) WriteFile(ctx context.Context, path string, data []byte) error {
	err := os.WriteFile(path, data, filePerm)
	if err == nil {
		return nil
	}
	if !errors.Is(err, fs.ErrPermission) && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return w.writeElevated(ctx, path, data)
}

func (w ElevatedWriter) RemoveFile(ctx context.Context, path string) error {
	err := os.Remove(path)
	if err == nil || errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if !errors.Is(err, fs.ErrPermission) {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	return w.Runner.Stream(ctx, sysexec.Command{
		Name:       removeBinary,
		Args:       []string{removeForce, endOfOptions, path},
		Privileged: true,
	})
}

func (w ElevatedWriter) writeElevated(ctx context.Context, path string, data []byte) error {
	tmp, err := os.CreateTemp("", tempFilePattern)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	return w.Runner.Stream(ctx, sysexec.Command{
		Name:       installBinary,
		Args:       []string{installDirs, installMode, fmt.Sprintf(fileModeFormat, filePerm), tmp.Name(), path},
		Privileged: true,
	})
}
