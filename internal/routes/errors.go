package routes

import (
	"errors"
	"fmt"
)

var (
	ErrTunnelNotFound     = errors.New("tunnel not found")
	ErrCredentialsMissing = errors.New("credentials file not found")
)

func (e *RollbackError) Error() string {
	return fmt.Sprintf(rollbackFormat, e.Cause)
}

func (e *RollbackError) Unwrap() error {
	return e.Cause
}

func (e *DNSRollbackError) Error() string {
	return fmt.Sprintf(dnsRollbackFormat, e.Cause)
}

func (e *DNSRollbackError) Unwrap() error {
	return e.Cause
}
