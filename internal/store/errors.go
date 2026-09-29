package store

import "errors"

var ErrMissing = errors.New("config file not found, run `cftun init` first")
