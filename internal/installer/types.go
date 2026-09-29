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

type ReleaseSource interface {
	Latest(ctx context.Context) (Release, error)
}

type FamilySource interface {
	Family() (Family, error)
}

type Installer struct {
	Runner     sysexec.Runner
	Downloader Downloader
	Releases   ReleaseSource
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

type Asset struct {
	Name   string
	URL    string
	SHA256 string
}

type Release struct {
	Assets []Asset
}

type GitHubReleases struct {
	Client *http.Client
}

type releasePayload struct {
	Assets []assetPayload `json:"assets"`
}

type assetPayload struct {
	Name   string `json:"name"`
	URL    string `json:"browser_download_url"`
	Digest string `json:"digest"`
}
