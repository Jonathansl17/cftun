package cli

func (e *presentedError) Error() string {
	return e.text
}

func (e *presentedError) Unwrap() error {
	return e.cause
}
