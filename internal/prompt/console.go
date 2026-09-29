package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

func NewConsole(in io.Reader, out io.Writer, texts Texts) *Console {
	return &Console{In: bufio.NewReader(in), Out: out, Texts: texts}
}

func (c *Console) Ask(label string, validate func(string) error) (string, error) {
	for {
		fmt.Fprintf(c.Out, c.Texts.PromptFormat, label)
		answer, err := c.line()
		if err != nil {
			return "", err
		}
		if err := validate(answer); err != nil {
			fmt.Fprintln(c.Out, err)
			continue
		}
		return answer, nil
	}
}

func (c *Console) Confirm(label string) (bool, error) {
	fmt.Fprintf(c.Out, c.Texts.ConfirmFormat, label)
	answer, err := c.line()
	if err != nil {
		return false, err
	}
	return yesAnswers[strings.ToLower(answer)], nil
}

func (c *Console) line() (string, error) {
	text, err := c.In.ReadString(lineDelimiter)
	if errors.Is(err, io.EOF) && text == "" {
		return "", ErrAborted
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf(readInputFormat, err)
	}
	return strings.TrimSpace(text), nil
}
