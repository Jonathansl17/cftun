package dnsapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c Client) Configured() (bool, error) {
	token, err := c.Tokens.Load()
	if err != nil {
		return false, err
	}
	return token != "", nil
}

func (c Client) DeleteCNAME(ctx context.Context, hostname string) (int, error) {
	zone, err := c.zoneID(ctx, hostname)
	if err != nil {
		return 0, err
	}
	query := url.Values{queryType: {recordType}, queryName: {hostname}}
	records, err := c.ids(ctx, http.MethodGet, fmt.Sprintf(recordsPathFormat, zone, query.Encode()))
	if err != nil {
		return 0, err
	}
	for deleted, id := range records {
		if _, err := c.call(ctx, http.MethodDelete, fmt.Sprintf(recordPathFormat, zone, id)); err != nil {
			return deleted, err
		}
	}
	return len(records), nil
}
