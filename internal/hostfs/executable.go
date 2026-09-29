package hostfs

import (
	"fmt"
	"os"
	"path/filepath"
)

func (FS) Executable() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locate running executable: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolve executable %s: %w", path, err)
	}
	return resolved, nil
}
