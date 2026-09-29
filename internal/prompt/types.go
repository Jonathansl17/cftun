package prompt

import (
	"bufio"
	"io"
)

type Prompter interface {
	Ask(label string, validate func(string) error) (string, error)
	Select(label string, options []string) (int, error)
	Confirm(label string) (bool, error)
}

type Texts struct {
	PromptFormat      string
	OptionFormat      string
	ConfirmFormat     string
	Choice            string
	ChoiceRangeFormat string
}

type Console struct {
	In    *bufio.Reader
	Out   io.Writer
	Texts Texts
}
