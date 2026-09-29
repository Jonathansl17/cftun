package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func (i Installer) Family() (Family, error) {
	return i.Families.Family()
}

func (i Installer) Install(ctx context.Context) error {
	family, err := i.Families.Family()
	if err != nil {
		return err
	}
	p := planFor(family)
	url, err := p.assetURL(i.GoArch)
	if err != nil {
		return err
	}
	if url == "" {
		return i.Runner.Stream(ctx, privileged(p.install("")))
	}
	return i.downloadAndInstall(ctx, p, url)
}

func (i Installer) downloadAndInstall(ctx context.Context, p plan, url string) error {
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
