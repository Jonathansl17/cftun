package uikit

import (
	"errors"
	"strings"
)

func Describe(err error) string {
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		return describeAll(joined.Unwrap())
	}
	var display *DisplayError
	if errors.As(err, &display) {
		return display.Text
	}
	for _, match := range matchers {
		if text, ok := match(err); ok {
			return text
		}
	}
	return err.Error()
}

func Displayable(err error) error {
	if err == nil {
		return nil
	}
	return NewDisplayError(Describe(err), err)
}

func NewDisplayError(text string, cause error) error {
	return &DisplayError{Text: text, Cause: cause}
}

func describeAll(errs []error) string {
	texts := make([]string, len(errs))
	for i, err := range errs {
		texts[i] = Describe(err)
	}
	return strings.Join(texts, joinSeparator)
}
