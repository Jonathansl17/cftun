package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

const downloadDirPattern = "cftun-download-*"

// Downloader fetches a URL into a local file.
type Downloader interface {
	Download(ctx context.Context, url, dest string) error
}

// Installer installs or removes cloudflared using the family's native method.
type Installer struct {
	Runner     sysexec.Runner
	Downloader Downloader
	Family     Family
	GoArch     string
}

// Install installs cloudflared.
func (i Installer) Install(ctx context.Context) error {
	p := planFor(i.Family)
	url, err := p.assetURL(i.GoArch)
	if err != nil {
		return err
	}
	if url == "" {
		return i.Runner.Stream(ctx, privileged(p.install("")))
	}
	dir, err := os.MkdirTemp("", downloadDirPattern)
	if err != nil {
		return fmt.Errorf("create download dir: %w", err)
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, filepath.Base(url))
	if err := i.Downloader.Download(ctx, url, file); err != nil {
		return err
	}
	return i.Runner.Stream(ctx, privileged(p.install(file)))
}

// Uninstall removes cloudflared through its package manager when it owns the
// binary, and deletes the static binary otherwise.
func (i Installer) Uninstall(ctx context.Context) error {
	p := planFor(i.Family)
	if !i.ownedByPackage(ctx, p) {
		p = binaryPlan
	}
	return i.Runner.Stream(ctx, privileged(p.remove))
}

// ownedByPackage reports whether the package manager has cloudflared
// registered. A failing query is the "not installed" answer, not an error.
func (i Installer) ownedByPackage(ctx context.Context, p plan) bool {
	if p.query == nil {
		return false
	}
	_, err := i.Runner.Output(ctx, sysexec.Command{Name: p.query[0], Args: p.query[1:]})
	return err == nil
}
