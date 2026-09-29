package hostfs

import (
	"errors"
	"fmt"
	"os"
)

func (FS) EmptyTemp(pattern string) (Temp, error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return Temp{}, fmt.Errorf("create temp file: %w", err)
	}
	if err := file.Close(); err != nil {
		removeErr := os.Remove(file.Name())
		return Temp{}, fmt.Errorf("close temp file: %w", errors.Join(err, removeErr))
	}
	return Temp{Path: file.Name()}, nil
}

func (t Temp) Remove() error {
	if err := os.Remove(t.Path); err != nil {
		return fmt.Errorf("remove temp file: %w", err)
	}
	return nil
}
