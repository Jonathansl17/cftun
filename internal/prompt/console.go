// Package prompt asks the user for input on the terminal.
package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Jonathansl17/cloudfare/internal/msg"
)

// ErrAborted reports that input ended before an answer was given.
var ErrAborted = errors.New("input aborted")

// Prompter collects answers from the user.
type Prompter interface {
	Ask(label string, validate func(string) error) (string, error)
	Select(label string, options []string) (int, error)
	Confirm(label string) (bool, error)
}

// Console is a Prompter over line-based streams.
type Console struct {
	In  *bufio.Reader
	Out io.Writer
}

// NewConsole wraps in and out.
func NewConsole(in io.Reader, out io.Writer) *Console {
	return &Console{In: bufio.NewReader(in), Out: out}
}

// Ask repeats the question until validate accepts the answer.
func (c *Console) Ask(label string, validate func(string) error) (string, error) {
	for {
		fmt.Fprintf(c.Out, msg.PromptFormat, label)
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

// Select shows numbered options and returns the chosen index.
func (c *Console) Select(label string, options []string) (int, error) {
	fmt.Fprintln(c.Out, label)
	for i, option := range options {
		fmt.Fprintf(c.Out, msg.OptionFormat, i+1, option)
	}
	answer, err := c.Ask(msg.PromptChoice, func(s string) error {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > len(options) {
			return fmt.Errorf(msg.ErrChoiceRange, len(options))
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	n, _ := strconv.Atoi(answer) // validated above
	return n - 1, nil
}

// Confirm asks a yes/no question that defaults to no.
func (c *Console) Confirm(label string) (bool, error) {
	fmt.Fprintf(c.Out, msg.ConfirmFormat, label)
	answer, err := c.line()
	if err != nil {
		return false, err
	}
	switch strings.ToLower(answer) {
	case msg.AnswerYes, msg.AnswerYesLong:
		return true, nil
	default:
		return false, nil
	}
}

func (c *Console) line() (string, error) {
	text, err := c.In.ReadString('\n')
	if errors.Is(err, io.EOF) && text == "" {
		return "", ErrAborted
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read input: %w", err)
	}
	return strings.TrimSpace(text), nil
}
