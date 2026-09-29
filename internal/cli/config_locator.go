package cli

import (
	"os"

	"github.com/Jonathansl17/cftun/internal/paths"
)

func newConfigLocator() *configLocator {
	return &configLocator{Env: os.Getenv}
}

func (l *configLocator) ConfigPath() string {
	if l.Flag != "" {
		return l.Flag
	}
	if env := l.Env(configPathEnv); env != "" {
		return env
	}
	return paths.DefaultConfig
}
