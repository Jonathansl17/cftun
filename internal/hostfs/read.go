package hostfs

import (
	"fmt"
	"os"
)

func (FS) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(readFileFormat, path, err)
	}
	return data, nil
}
