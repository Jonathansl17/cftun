package dnsapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultBaseURL = "https://api.cloudflare.com/client/v4"
	TokenEnv       = "CLOUDFLARE_API_TOKEN"
	recordType     = "CNAME"
	minZoneLabels  = 2
)

var ErrZoneNotFound = errors.New("no Cloudflare zone found for hostname")

type Client struct {
	HTTP    *http.Client
	BaseURL string
	Token   string
}

type envelope struct {
	Success bool              `json:"success"`
	Errors  []json.RawMessage `json:"errors"`
	Result  []struct {
		ID string `json:"id"`
	} `json:"result"`
}

func (c Client) Configured() bool {
	return c.Token != ""
}

func (c Client) DeleteCNAME(ctx context.Context, hostname string) (int, error) {
	zone, err := c.zoneID(ctx, hostname)
	if err != nil {
		return 0, err
	}
	query := url.Values{"type": {recordType}, "name": {hostname}}
	records, err := c.ids(ctx, http.MethodGet, "/zones/"+zone+"/dns_records?"+query.Encode())
	if err != nil {
		return 0, err
	}
	for _, id := range records {
		if _, err := c.call(ctx, http.MethodDelete, "/zones/"+zone+"/dns_records/"+id); err != nil {
			return 0, err
		}
	}
	return len(records), nil
}

func (c Client) zoneID(ctx context.Context, hostname string) (string, error) {
	labels := strings.Split(hostname, ".")
	for i := 0; i <= len(labels)-minZoneLabels; i++ {
		name := strings.Join(labels[i:], ".")
		ids, err := c.ids(ctx, http.MethodGet, "/zones?"+url.Values{"name": {name}}.Encode())
		if err != nil {
			return "", err
		}
		if len(ids) > 0 {
			return ids[0], nil
		}
	}
	return "", fmt.Errorf("%s: %w", hostname, ErrZoneNotFound)
}

func (c Client) ids(ctx context.Context, method, path string) ([]string, error) {
	env, err := c.call(ctx, method, path)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(env.Result))
	for _, r := range env.Result {
		ids = append(ids, r.ID)
	}
	return ids, nil
}

func (c Client) call(ctx context.Context, method, path string) (*envelope, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cloudflare api: %w", err)
	}
	defer resp.Body.Close()
	var env envelope
	if method == http.MethodDelete {
		return &env, statusError(resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("cloudflare api: decode: %w", err)
	}
	if !env.Success {
		return nil, fmt.Errorf("cloudflare api %s: %s", resp.Status, env.Errors)
	}
	return &env, nil
}

func statusError(resp *http.Response) error {
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("cloudflare api: unexpected status %s", resp.Status)
	}
	return nil
}
