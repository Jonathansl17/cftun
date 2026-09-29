package installer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (g GitHubReleases) Latest(ctx context.Context) (Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestReleaseURL, nil)
	if err != nil {
		return Release{}, fmt.Errorf(releaseRequestFormat, err)
	}
	req.Header.Set(acceptHeader, acceptGitHubJSON)
	req.Header.Set(apiVersionHeader, apiVersion)
	resp, err := g.Client.Do(req)
	if err != nil {
		return Release{}, fmt.Errorf(releaseFetchFormat, err)
	}
	release, err := decodeRelease(resp)
	if closeErr := resp.Body.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf(releaseFetchFormat, closeErr)
	}
	return release, err
}

func decodeRelease(resp *http.Response) (Release, error) {
	if resp.StatusCode != http.StatusOK {
		return Release{}, fmt.Errorf(releaseStatusFormat, resp.Status)
	}
	var payload releasePayload
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxReleaseBytes)).Decode(&payload); err != nil {
		return Release{}, fmt.Errorf(releaseDecodeFormat, err)
	}
	assets := make([]Asset, 0, len(payload.Assets))
	for _, a := range payload.Assets {
		assets = append(assets, Asset{Name: a.Name, URL: a.URL, SHA256: sha256Of(a.Digest)})
	}
	return Release{Assets: assets}, nil
}

func sha256Of(digest string) string {
	hexSum, ok := strings.CutPrefix(digest, digestPrefix)
	if !ok {
		return ""
	}
	return hexSum
}

func (r Release) Find(name string) (Asset, error) {
	for _, a := range r.Assets {
		if a.Name == name {
			return a, nil
		}
	}
	return Asset{}, fmt.Errorf(assetNotFoundFormat, ErrAssetNotFound, name)
}
