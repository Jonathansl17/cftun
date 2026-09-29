package health

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"syscall"

	"github.com/Jonathansl17/cftun/internal/msg"
)

const (
	maxBodyBytes     = 4096
	statusTunnelDown = 530
	errorCodeDNS     = "1016"
	publicScheme     = "https://"
)

type Result struct {
	URL    string
	Status int
	Err    error
	dns    bool
}

func (r Result) OK() bool {
	return r.Err == nil && r.Status < http.StatusInternalServerError
}

func (r Result) Hint() string {
	switch {
	case r.OK():
		return ""
	case errors.Is(r.Err, syscall.ECONNREFUSED):
		return msg.HintConnectionRefused
	case r.dns:
		return msg.HintDNS1016
	case r.Status == http.StatusBadGateway || r.Status == statusTunnelDown:
		return msg.HintTunnelDown
	default:
		return msg.HintGeneric
	}
}

type Checker struct {
	HTTP *http.Client
}

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
