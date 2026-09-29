package dnsapi

import (
	"errors"
	"fmt"
)

var (
	ErrZoneNotFound = errors.New("no Cloudflare zone found for hostname")
	ErrNoToken      = errors.New("no Cloudflare API token")
)

func (e *APIError) Error() string {
	return fmt.Sprintf(apiErrorFormat, e.Op, e.Err)
}

func (e *APIError) Unwrap() error {
	return e.Err
}
