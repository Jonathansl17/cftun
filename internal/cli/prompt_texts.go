package cli

import (
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/prompt"
)

func promptTexts() prompt.Texts {
	return prompt.Texts{
		PromptFormat:      msg.PromptFormat,
		OptionFormat:      msg.OptionFormat,
		ConfirmFormat:     msg.ConfirmFormat,
		Choice:            msg.PromptChoice,
		ChoiceRangeFormat: msg.ErrChoiceRange,
	}
}
