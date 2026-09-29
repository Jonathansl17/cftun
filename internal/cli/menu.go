package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/prompt"
)

func runMenu(a *App, root *cobra.Command) error {
	return browse(a, root, menuScreen{title: msg.MenuTitle, entries: mainMenu, leave: msg.MenuExit})
}

func browse(a *App, root *cobra.Command, screen menuScreen) error {
	for {
		labels := make([]string, 0, len(screen.entries)+1)
		for _, e := range screen.entries {
			labels = append(labels, e.label)
		}
		i, err := a.Prompt.Select(screen.title, append(labels, screen.leave))
		if errors.Is(err, prompt.ErrAborted) || (err == nil && i == len(screen.entries)) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := dispatch(a, root, screen.entries[i]); err != nil {
			return err
		}
		if a.CftunRemoved {
			return nil
		}
	}
}

func dispatch(a *App, root *cobra.Command, e entry) error {
	if e.children != nil {
		return browse(a, root, menuScreen{title: e.label, entries: e.children, leave: msg.MenuBack})
	}
	cmd, _, err := root.Find(e.path)
	if err != nil {
		return err
	}
	cmd.SetContext(root.Context())
	err = a.Gate.Require(cmd)
	if err == nil {
		err = cmd.RunE(cmd, nil)
	}
	if err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			return err
		}
		a.PrintError(err)
	}
	return nil
}
