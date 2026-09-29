package uikit

import (
	"fmt"

	"github.com/Jonathansl17/cftun/internal/msg"
)

func (s Session) Printf(format string, args ...any) {
	fmt.Fprintf(s.Out, format, args...)
}

func (s Session) PrintError(err error) {
	s.Printf(msg.ErrorFormat, Describe(err))
}

func (s Session) ValueOrAsk(value, label string, validate func(string) error) (string, error) {
	if value != "" {
		return value, validate(value)
	}
	return s.Prompt.Ask(label, validate)
}

func (s Session) ConfirmOrCancel(label string) error {
	ok, err := s.Prompt.Confirm(label)
	if err != nil {
		return err
	}
	if !ok {
		return ErrCancelled
	}
	return nil
}
