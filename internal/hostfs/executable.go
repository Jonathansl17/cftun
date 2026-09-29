package hostfs

import (
	"fmt"
	"os"
	"path/filepath"
)

func (FS) Executable() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf(locateExecutableFormat, err)
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf(resolveExecutableFormat, path, err)
	}
	return resolved, nil
}
