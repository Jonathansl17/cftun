package teardown

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/hostfs"
	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func (r SafeRemover) RemoveAll(ctx context.Context, targets []Target) error {
	recursive, plain, err := r.plan(targets)
	if err != nil {
		return err
	}
	return errors.Join(r.remove(ctx, recursiveFlags, recursive), r.remove(ctx, plainFlags, plain))
}

func (r SafeRemover) plan(targets []Target) ([]string, []string, error) {
	var recursive, plain []string
	for _, target := range targets {
		remove, err := r.vet(target)
		if err != nil {
			return nil, nil, err
		}
		switch {
		case !remove:
		case target.Kind == TargetDirectory:
			recursive = append(recursive, target.Path)
		default:
			plain = append(plain, target.Path)
		}
	}
	return recursive, plain, nil
}

func (r SafeRemover) vet(target Target) (bool, error) {
	if !filepath.IsAbs(target.Path) || filepath.Clean(target.Path) != target.Path {
		return false, &UnsafeTargetError{Path: target.Path, Kind: target.Kind}
	}
	if isSystemTarget(target) {
		return true, nil
	}
	kind, err := r.Inspector.KindOf(target.Path)
	if err != nil {
		return false, err
	}
	if kind == hostfs.KindMissing {
		return false, nil
	}
	if !allowedUserTarget(target, kind) {
		return false, &UnsafeTargetError{Path: target.Path, Kind: target.Kind}
	}
	return true, nil
}

func (r SafeRemover) remove(ctx context.Context, flags string, targets []string) error {
	if len(targets) == 0 {
		return nil
	}
	args := append([]string{flags, endOfOptions}, targets...)
	return r.Runner.Stream(ctx, sysexec.Command{Name: removeBinary, Args: args, Privileged: true})
}

func isSystemTarget(target Target) bool {
	for _, known := range systemTargets {
		if known == target {
			return true
		}
	}
	return false
}

func allowedUserTarget(target Target, kind hostfs.Kind) bool {
	switch target.Kind {
	case TargetFile:
		return kind == hostfs.KindRegular
	case TargetDirectory:
		return kind == hostfs.KindDirectory && filepath.Base(target.Path) == paths.UserDirName
	default:
		return false
	}
}
