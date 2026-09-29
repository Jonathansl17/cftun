package uikit

import (
	"io"

	"github.com/Jonathansl17/cftun/internal/prompt"
)

type Session struct {
	Out    io.Writer
	Prompt prompt.Prompter
}

type Matcher func(err error) (string, bool)

type DisplayError struct {
	Text  string
	Cause error
}
