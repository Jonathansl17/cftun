package uikit

import (
	"errors"

	"github.com/Jonathansl17/cftun/internal/msg"
)

var (
	ErrCancelled          = errors.New(msg.ErrCancelled)
	ErrEmpty              = errors.New(msg.ErrEmpty)
	ErrCloudflaredMissing = errors.New(msg.ErrCloudflaredMissing)
)

func (e *DisplayError) Error() string {
	return e.Text
}

func (e *DisplayError) Unwrap() error {
	return e.Cause
}
