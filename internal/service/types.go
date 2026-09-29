package service

import "github.com/Jonathansl17/cftun/internal/sysexec"

type Action string

type Controller struct {
	Name     string
	Runner   sysexec.Runner
	commands map[Action][]string
	exists   func(string) bool
}

type UnsupportedActionError struct {
	Name   string
	Action Action
}
