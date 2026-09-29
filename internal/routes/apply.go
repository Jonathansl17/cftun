package routes

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jonathansl17/cftun/internal/ingress"
)

func (a applier) apply(ctx context.Context, mutate func(*ingress.Document) error) (*ingress.Document, string, error) {
	doc, err := a.store.Load()
	if err != nil {
		return nil, "", err
	}
	if err := mutate(doc); err != nil {
		return nil, "", err
	}
	validation, err := a.commit(ctx, doc, a.restore)
	return doc, validation, err
}

func (a applier) commit(ctx context.Context, doc *ingress.Document, rollback func(context.Context) error) (string, error) {
	if err := a.store.Save(ctx, doc); err != nil {
		return "", err
	}
	validation, err := a.checker.Validate(ctx)
	if err == nil {
		return validation, nil
	}
	if rerr := rollback(ctx); rerr != nil {
		return "", errors.Join(err, rerr)
	}
	return "", &RollbackError{Cause: err}
}

func (a applier) undo(ctx context.Context, cause error) error {
	if err := a.restore(ctx); err != nil {
		return errors.Join(cause, err)
	}
	return &DNSRollbackError{Cause: cause}
}

func (a applier) restore(ctx context.Context) error {
	if err := a.store.Restore(ctx); err != nil {
		return fmt.Errorf(restoreBackupFormat, err)
	}
	return nil
}

func (a applier) discard(ctx context.Context) error {
	if err := a.store.Delete(ctx); err != nil {
		return fmt.Errorf(removeNewConfigFormat, err)
	}
	return nil
}
