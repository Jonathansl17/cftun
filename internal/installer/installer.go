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
	name, err := p.assetName(i.GoArch)
	if err != nil {
		return err
	}
	if name == "" {
		return i.Runner.Stream(ctx, privileged(p.install("")))
	}
	return i.downloadAndInstall(ctx, p, name)
}

func (i Installer) downloadAndInstall(ctx context.Context, p plan, name string) error {
	asset, err := i.findAsset(ctx, name)
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", downloadDirPattern)
	if err != nil {
		return fmt.Errorf(createDirFormat, err)
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, asset.Name)
	if err := i.Downloader.Download(ctx, asset.URL, file); err != nil {
		return err
	}
	if err := VerifySHA256(file, asset.SHA256); err != nil {
		return err
	}
	return i.Runner.Stream(ctx, privileged(p.install(file)))
}

func (i Installer) findAsset(ctx context.Context, name string) (Asset, error) {
	release, err := i.Releases.Latest(ctx)
	if err != nil {
		return Asset{}, err
	}
	asset, err := release.Find(name)
	if err != nil {
		return Asset{}, err
	}
	if asset.SHA256 == "" {
		return Asset{}, fmt.Errorf(missingDigestFormat, ErrMissingDigest, name)
	}
	return asset, nil
}
