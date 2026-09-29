package installcmd

import (
	"fmt"
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/teardown"
)

func (g *Group) teardownParams() (teardown.Params, error) {
	config, err := filepath.Abs(g.deps.Config.Path())
	if err != nil {
		return teardown.Params{}, fmt.Errorf(resolveFormat, err)
	}
	backup, err := filepath.Abs(g.deps.Config.BackupPath())
	if err != nil {
		return teardown.Params{}, fmt.Errorf(resolveFormat, err)
	}
	userDir, err := uikit.UserCloudflaredDir(g.deps.Home.Home)
	if err != nil {
		return teardown.Params{}, err
	}
	return teardown.Params{ConfigPath: config, BackupPath: backup, UserDir: userDir}, nil
}
