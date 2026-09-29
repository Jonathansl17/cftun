package hostfs

import (
	"fmt"
	"os"
)

func (FS) Home() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("locate home: %w", err)
	}
	return home, nil
}
