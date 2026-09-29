package installer

import (
	"context"
	"fmt"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func (i Installer) Uninstall(ctx context.Context) error {
	family, err := i.Families.Family()
	if err != nil {
		return err
	}
	p := planFor(family)
	owned, err := i.ownedByPackage(ctx, p)
	if err != nil {
		return err
	}
	if !owned {
		p = binaryPlan
	}
	return i.Runner.Stream(ctx, privileged(p.remove))
}

func (i Installer) ownedByPackage(ctx context.Context, p plan) (bool, error) {
	if p.query == nil {
		return false, nil
	}
	_, err := i.Runner.Output(ctx, sysexec.Command{Name: p.query[0], Args: p.query[1:]})
	if err == nil {
		return true, nil
	}
	if code, ok := sysexec.ExitCode(err); ok && code == queryNotInstalledExitCode {
		return false, nil
	}
	return false, fmt.Errorf(queryFailedFormat, p.query[0], err)
}
