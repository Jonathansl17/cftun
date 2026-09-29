package uikit

import (
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/paths"
)

func UserCloudflaredDir(home func() (string, error)) (string, error) {
	dir, err := home()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, paths.UserDirName), nil
}
