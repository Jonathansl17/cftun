// Package store loads and saves the cloudflared config file with backups.
package store

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

const (
	filePerm        = 0o644
	filePermOctal   = "0644"
	installBinary   = "install"
	tempFilePattern = "cftun-*"
)

// Writer persists bytes to a path.
type Writer interface {
	WriteFile(ctx context.Context, path string, data []byte) error
}

// ElevatedWriter writes directly and, when the target is not writable by the
// current user, falls back to a privileged install of a temporary copy.
type ElevatedWriter struct {
	Runner sysexec.Runner
}

// WriteFile implements Writer.
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
		Args:       []string{"-D", "-m", filePermOctal, tmp.Name(), path},
		Privileged: true,
	})
}
