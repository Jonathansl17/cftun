package uikit

import (
	"io"

	"github.com/Jonathansl17/cftun/internal/prompt"
)

type Session struct {
	Out    io.Writer
	Prompt prompt.Prompter
}

type Matcher func(err error) (string, bool)

type DisplayError struct {
	Text  string
	Cause error
}

type MenuGroup string

const (
	GroupNone    MenuGroup = ""
	GroupRoutes  MenuGroup = "routes"
	GroupService MenuGroup = "service"
	GroupTunnels MenuGroup = "tunnels"
	GroupAccount MenuGroup = "account"
)

type MenuSlot struct {
	Group MenuGroup
	Label string
	Order int
}
