package cli

import (
	"context"
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/msg"
)

func removeCftun(ctx context.Context, a *App) error {
	binary, err := a.Files.Executable()
	if err != nil {
		return err
	}
	targets := []string{binary}
	for _, file := range completionFiles {
		if a.Files.Exists(file) {
			targets = append(targets, file)
		}
	}
	for _, target := range targets {
		a.Printf(msg.InfoRemoving, filepath.Clean(target))
		if err := a.Store.Writer.RemoveFile(ctx, target); err != nil {
			return err
		}
	}
	a.Printf(msg.InfoCftunRemoved)
	a.CftunRemoved = true
	return nil
}
