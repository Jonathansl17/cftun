package installer

import (
	"context"
	"net/http"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

type Family int

type Downloader interface {
	Download(ctx context.Context, url, dest string) error
}

type FamilySource interface {
	Family() (Family, error)
}

type Installer struct {
	Runner     sysexec.Runner
	Downloader Downloader
	Families   FamilySource
	GoArch     string
}

type plan struct {
	assetFormat   string
	archOverrides map[string]string
	install       func(file string) []string
	query         []string
	remove        []string
}

type HTTPDownloader struct {
	Client *http.Client
}

type OSReleaseSource struct{}
