package health

import (
	"context"
	"io"
	"net/http"
	"strings"
)

func PublicURL(hostname string) string {
	return publicScheme + hostname
}

func (c Checker) Probe(ctx context.Context, url string) Result {
	res := Result{URL: url}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		res.Err = err
		return res
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		res.Err = err
		return res
	}
	defer resp.Body.Close()
	res.Status = resp.StatusCode
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		res.Err = err
		return res
	}
	res.dns = res.Status == statusTunnelDown && strings.Contains(string(body), errorCodeDNS)
	return res
}
