package installer

import (
	"fmt"
	"slices"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func planFor(f Family) plan {
	if p, ok := plans[f]; ok {
		return p
	}
	return binaryPlan
}

func (p plan) assetURL(goarch string) (string, error) {
	if p.assetFormat == "" {
		return "", nil
	}
	if !slices.Contains(supportedArches, goarch) {
		return "", fmt.Errorf(unsupportedArchFormat, ErrUnsupportedArch, goarch)
	}
	return releaseBaseURL + fmt.Sprintf(p.assetFormat, p.archName(goarch)), nil
}

func (p plan) archName(goarch string) string {
	if name, ok := p.archOverrides[goarch]; ok {
		return name
	}
	return goarch
}

func privileged(args []string) sysexec.Command {
	return sysexec.Command{Name: args[0], Args: args[1:], Privileged: true}
}
