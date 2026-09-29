package installcmd

import (
	"context"
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) removeSelf(ctx context.Context) error {
	binary, err := g.deps.Files.Executable()
	if err != nil {
		return err
	}
	targets := []string{binary}
	for _, file := range completionFiles {
		if g.deps.Files.Exists(file) {
			targets = append(targets, file)
		}
	}
	for _, target := range targets {
		g.deps.Session.Printf(msg.InfoRemoving, filepath.Clean(target))
		if err := g.deps.Remover.RemoveFile(ctx, target); err != nil {
			return err
		}
	}
	g.deps.Session.Printf(msg.InfoCftunRemoved)
	g.removed = true
	return nil
}
