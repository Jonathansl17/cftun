package store

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Jonathansl17/cftun/internal/ingress"
)

const backupSuffix = ".bak"

var ErrMissing = errors.New("config file not found, run `cftun init` first")

type File struct {
	Path   string
	Writer Writer
}

func (f File) Exists() bool {
	_, err := os.Stat(f.Path)
	return err == nil
}

func (f File) BackupPath() string {
	return f.Path + backupSuffix
}

func (f File) Load() (*ingress.Document, error) {
	data, err := os.ReadFile(f.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%s: %w", f.Path, ErrMissing)
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", f.Path, err)
	}
	return ingress.Parse(data)
}

func (f File) Save(ctx context.Context, doc *ingress.Document) error {
	data, err := doc.Bytes()
	if err != nil {
		return err
	}
	if err := f.backup(ctx); err != nil {
		return err
	}
	return f.Writer.WriteFile(ctx, f.Path, data)
}

func (f File) Restore(ctx context.Context) error {
	data, err := os.ReadFile(f.BackupPath())
	if err != nil {
		return fmt.Errorf("read backup: %w", err)
	}
	return f.Writer.WriteFile(ctx, f.Path, data)
}

func (f File) backup(ctx context.Context) error {
	data, err := os.ReadFile(f.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", f.Path, err)
	}
	return f.Writer.WriteFile(ctx, f.BackupPath(), data)
}
