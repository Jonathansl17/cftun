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
		mode, err := r.vet(target)
		if err != nil {
			return nil, nil, err
		}
		switch mode {
		case removalRecursive:
			recursive = append(recursive, target.Path)
		case removalPlain:
			plain = append(plain, target.Path)
		}
	}
	return recursive, plain, nil
}

func (r SafeRemover) vet(target Target) (removal, error) {
	if !filepath.IsAbs(target.Path) || filepath.Clean(target.Path) != target.Path {
		return removalSkip, &UnsafeTargetError{Path: target.Path, Kind: target.Kind}
	}
	if isSystemTarget(target) {
		return kindRemoval(target), nil
	}
	kind, err := r.Inspector.KindOf(target.Path)
	if err != nil {
		return removalSkip, err
	}
	switch {
	case kind == hostfs.KindMissing:
		return removalSkip, nil
	case kind == hostfs.KindSymlink:
		return removalPlain, nil
	case !allowedUserTarget(target, kind):
		return removalSkip, &UnsafeTargetError{Path: target.Path, Kind: target.Kind}
	}
	return kindRemoval(target), nil
}

func kindRemoval(target Target) removal {
	if target.Kind == TargetDirectory {
		return removalRecursive
	}
	return removalPlain
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
