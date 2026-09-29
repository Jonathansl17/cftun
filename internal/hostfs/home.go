package hostfs

import (
	"fmt"
	"os"
)

func (FS) Home() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf(locateHomeFormat, err)
	}
	return home, nil
}
