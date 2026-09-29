package store

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

type Locator interface {
	ConfigPath() string
}

type Writer interface {
	WriteFile(ctx context.Context, path string, data []byte) error
	RemoveFile(ctx context.Context, path string) error
}

type File struct {
	Locator Locator
	Writer  Writer
}

type ElevatedWriter struct {
	Runner sysexec.Runner
}
