package dnsapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c Client) zoneID(ctx context.Context, hostname string) (string, error) {
	labels := strings.Split(hostname, labelSeparator)
	for i := 0; i <= len(labels)-minZoneLabels; i++ {
		name := strings.Join(labels[i:], labelSeparator)
		query := url.Values{queryName: {name}}
		ids, err := c.ids(ctx, http.MethodGet, fmt.Sprintf(zonesPathFormat, query.Encode()))
		if err != nil {
			return "", err
		}
		if len(ids) > 0 {
			return ids[0], nil
		}
	}
	return "", fmt.Errorf(zoneErrorFormat, hostname, ErrZoneNotFound)
}
