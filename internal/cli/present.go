package cli

import (
	"errors"
	"fmt"

	"github.com/Jonathansl17/cftun/internal/apperr"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/routes"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

func present(err error) error {
	var rollback *routes.RollbackError
	var dnsRollback *routes.DNSRollbackError
	var unsafe *teardown.UnsafeTargetError
	switch {
	case err == nil:
		return nil
	case errors.As(err, &rollback):
		return &presentedError{text: fmt.Sprintf(causeTextFormat, msg.ErrRolledBack, rollback.Cause), cause: err}
	case errors.As(err, &dnsRollback):
		return &presentedError{text: fmt.Sprintf(causeTextFormat, msg.ErrDNSRolledBack, dnsRollback.Cause), cause: err}
	case errors.As(err, &unsafe):
		return &presentedError{text: fmt.Sprintf(msg.ErrUnsafeTarget, unsafe.Path), cause: err}
	}
	return presentSentinel(err)
}

func presentSentinel(err error) error {
	for sentinel, text := range sentinelTexts {
		if !errors.Is(err, sentinel) {
			continue
		}
		if subject, ok := apperr.SubjectOf(err); ok {
			text = fmt.Sprintf(subjectTextFormat, subject, text)
		}
		return &presentedError{text: text, cause: err}
	}
	return err
}
