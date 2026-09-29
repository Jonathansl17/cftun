package hostfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (FS) EmptyTemp(pattern string) (Temp, error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return Temp{}, fmt.Errorf(createTempFormat, err)
	}
	if err := file.Close(); err != nil {
		removeErr := os.Remove(file.Name())
		return Temp{}, fmt.Errorf(closeTempFormat, errors.Join(err, removeErr))
	}
	return Temp{Path: file.Name()}, nil
}

func (t Temp) Remove() error {
	if err := os.Remove(t.Path); err != nil {
		return fmt.Errorf(removeTempFormat, err)
	}
	return nil
}

func (FS) TempMatches(pattern string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(os.TempDir(), pattern))
	if err != nil {
		return nil, fmt.Errorf(globTempFormat, err)
	}
	return matches, nil
}

func (t Temp) Write(content string) error {
	if err := os.WriteFile(t.Path, []byte(content), tempFilePerm); err != nil {
		return fmt.Errorf(writeTempFormat, err)
	}
	return nil
}
