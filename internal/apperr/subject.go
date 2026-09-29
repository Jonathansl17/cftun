package apperr

import "errors"

func Wrap(subject string, err error) error {
	return &Subject{Name: subject, Err: err}
}

func SubjectOf(err error) (string, bool) {
	var s *Subject
	if errors.As(err, &s) {
		return s.Name, true
	}
	return "", false
}

func (s *Subject) Error() string {
	return s.Name + subjectSeparator + s.Err.Error()
}

func (s *Subject) Unwrap() error {
	return s.Err
}
