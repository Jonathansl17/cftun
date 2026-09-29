package store

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/Jonathansl17/cftun/internal/apperr"
	"github.com/Jonathansl17/cftun/internal/ingress"
)

func (f File) Path() string {
	return f.Locator.ConfigPath()
}

func (f File) BackupPath() string {
	return f.Path() + backupSuffix
}

func (f File) Exists() bool {
	return f.Reader.Exists(f.Path())
}

func (f File) Load() (*ingress.Document, error) {
	data, err := f.Reader.ReadFile(f.Path())
	if errors.Is(err, fs.ErrNotExist) {
		return nil, apperr.Wrap(f.Path(), ErrMissing)
	}
	if err != nil {
		return nil, err
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
	return f.Writer.WriteFile(ctx, f.Path(), data)
}

func (f File) Restore(ctx context.Context) error {
	data, err := f.Reader.ReadFile(f.BackupPath())
	if err != nil {
		return fmt.Errorf(readBackupFormat, err)
	}
	return f.Writer.WriteFile(ctx, f.Path(), data)
}

func (f File) Delete(ctx context.Context) error {
	return f.Writer.RemoveFile(ctx, f.Path())
}

func (f File) backup(ctx context.Context) error {
	data, err := f.Reader.ReadFile(f.Path())
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return f.Writer.WriteFile(ctx, f.BackupPath(), data)
}
