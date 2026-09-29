package prompt

import (
	"fmt"
	"strconv"
)

func (c *Console) Select(label string, options []string) (int, error) {
	fmt.Fprintln(c.Out, label)
	for i, option := range options {
		fmt.Fprintf(c.Out, c.Texts.OptionFormat, i+1, option)
	}
	chosen := 0
	_, err := c.Ask(c.Texts.Choice, func(answer string) error {
		n, err := strconv.Atoi(answer)
		if err != nil || n < 1 || n > len(options) {
			return fmt.Errorf(c.Texts.ChoiceRangeFormat, len(options))
		}
		chosen = n - 1
		return nil
	})
	if err != nil {
		return 0, err
	}
	return chosen, nil
}
