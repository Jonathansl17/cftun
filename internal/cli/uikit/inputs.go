package uikit

import "github.com/Jonathansl17/cftun/internal/ingress"

func AskHostname(s Session, value, label string) (string, error) {
	return AskParsed(s, value, label, ingress.NormalizeHostname)
}

func AskPort(s Session, value, label string) (int, error) {
	return AskParsed(s, value, label, ingress.ParsePort)
}

func NotEmpty(value string) error {
	if value == "" {
		return ErrEmpty
	}
	return nil
}

func FirstArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}
