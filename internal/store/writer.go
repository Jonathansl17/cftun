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
		return fmt.Errorf(writeFileFormat, path, err)
	}
	return w.writeElevated(ctx, path, data)
}

func (w ElevatedWriter) RemoveFile(ctx context.Context, path string) error {
	err := os.Remove(path)
	if err == nil || errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if !errors.Is(err, fs.ErrPermission) {
		return fmt.Errorf(removeFileFormat, path, err)
	}
	return w.Runner.Stream(ctx, sysexec.Command{
		Name:       removeBinary,
		Args:       []string{removeForce, endOfOptions, path},
		Privileged: true,
	})
}

func (w ElevatedWriter) writeElevated(ctx context.Context, path string, data []byte) (err error) {
	tmp, err := os.CreateTemp("", tempFilePattern)
	if err != nil {
		return fmt.Errorf(createTempFormat, err)
	}
	defer func() {
		if removeErr := os.Remove(tmp.Name()); removeErr != nil {
			err = errors.Join(err, fmt.Errorf(removeTempFormat, removeErr))
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		closeErr := tmp.Close()
		return fmt.Errorf(writeTempFormat, errors.Join(err, closeErr))
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf(closeTempFormat, err)
	}
	return w.Runner.Stream(ctx, sysexec.Command{
		Name:       installBinary,
		Args:       []string{installDirs, installMode, fmt.Sprintf(fileModeFormat, filePerm), tmp.Name(), path},
		Privileged: true,
	})
}
