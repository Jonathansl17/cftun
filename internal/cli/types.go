package cli

type Reporter interface {
	Printf(format string, args ...any)
}

type teardownPrinter struct {
	Report Reporter
}

type presentedError struct {
	text  string
	cause error
}
