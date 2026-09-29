package installer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (d HTTPDownloader) Download(ctx context.Context, rawURL, dest string) error {
	if err := requireHTTPS(rawURL); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf(buildRequestFormat, err)
	}
	client := *d.Client
	client.CheckRedirect = httpsOnlyRedirect
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(downloadFormat, rawURL, err)
	}
	err = saveBody(resp, rawURL, dest)
	if closeErr := resp.Body.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf(downloadFormat, rawURL, closeErr)
	}
	return err
}

func saveBody(resp *http.Response, rawURL, dest string) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(downloadStatusFormat, rawURL, resp.Status)
	}
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, downloadPerm)
	if err != nil {
		return fmt.Errorf(createFileFormat, dest, err)
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		return errors.Join(fmt.Errorf(saveFileFormat, dest, err), out.Close())
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf(saveFileFormat, dest, err)
	}
	return nil
}

func requireHTTPS(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf(parseURLFormat, rawURL, err)
	}
	if u.Scheme != httpsScheme {
		return fmt.Errorf(insecureURLFormat, ErrInsecureURL, rawURL)
	}
	return nil
}

func httpsOnlyRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return http.ErrUseLastResponse
	}
	if req.URL.Scheme != httpsScheme {
		return fmt.Errorf(insecureRedirectFormat, ErrInsecureRedirect, req.URL)
	}
	return nil
}
