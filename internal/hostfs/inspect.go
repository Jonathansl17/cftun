package hostfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func (FS) Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

func (FS) KindOf(path string) (Kind, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return KindMissing, nil
	}
	if err != nil {
		return KindOther, fmt.Errorf(inspectFormat, path, err)
	}
	switch mode := info.Mode(); {
	case mode&os.ModeSymlink != 0:
		return KindSymlink, nil
	case mode.IsRegular():
		return KindRegular, nil
	case mode.IsDir():
		return KindDirectory, nil
	default:
		return KindOther, nil
	}
}
