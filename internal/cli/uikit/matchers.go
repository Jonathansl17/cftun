package uikit

import (
	"errors"
	"fmt"

	"github.com/Jonathansl17/cftun/internal/apperr"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/routes"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

func sentinel(target error, text string) Matcher {
	return func(err error) (string, bool) {
		if !errors.Is(err, target) {
			return "", false
		}
		if subject, ok := apperr.SubjectOf(err); ok {
			return fmt.Sprintf(subjectTextFormat, subject, text), true
		}
		return text, true
	}
}

func invalidPortText() string {
	return fmt.Sprintf(msg.ErrInvalidPortFormat, ingress.MinPort, ingress.MaxPort)
}

func matchRollback(err error) (string, bool) {
	var rollback *routes.RollbackError
	if !errors.As(err, &rollback) {
		return "", false
	}
	return fmt.Sprintf(causeTextFormat, msg.ErrRolledBack, rollback.Cause), true
}

func matchDNSRollback(err error) (string, bool) {
	var rollback *routes.DNSRollbackError
	if !errors.As(err, &rollback) {
		return "", false
	}
	return fmt.Sprintf(causeTextFormat, msg.ErrDNSRolledBack, rollback.Cause), true
}

func matchUnsafeTarget(err error) (string, bool) {
	var unsafe *teardown.UnsafeTargetError
	if !errors.As(err, &unsafe) {
		return "", false
	}
	return fmt.Sprintf(msg.ErrUnsafeTarget, unsafe.Path), true
}
