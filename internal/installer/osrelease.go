package installer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Jonathansl17/cftun/internal/paths"
)

func (OSReleaseSource) Family() (Family, error) {
	f, err := os.Open(paths.OSRelease)
	if errors.Is(err, fs.ErrNotExist) {
		return FamilyUnknown, nil
	}
	if err != nil {
		return FamilyUnknown, fmt.Errorf(readOSReleaseFormat, paths.OSRelease, err)
	}
	defer f.Close()
	return DetectFamily(f)
}
