package ingress

import "errors"

var (
	ErrMalformed       = errors.New("malformed ingress config")
	ErrDuplicate       = errors.New("duplicate hostname")
	ErrNotFound        = errors.New("hostname not found")
	ErrInvalidHostname = errors.New("invalid hostname")
	ErrInvalidPort     = errors.New("invalid port")
)
