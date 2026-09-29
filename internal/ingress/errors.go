package ingress

import (
	"errors"
	"fmt"
)

var (
	ErrMalformed       = errors.New("config must be a mapping with an ingress list")
	ErrDuplicate       = errors.New("hostname already has a rule")
	ErrNotFound        = errors.New("hostname has no rule")
	ErrInvalidHostname = errors.New("invalid hostname, expected something like app.example.com")
	ErrInvalidPort     = fmt.Errorf("invalid port, expected a number between %d and %d", MinPort, MaxPort)
)
