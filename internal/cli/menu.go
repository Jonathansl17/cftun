package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/prompt"
)

// entry is one menu line: either a command path or a submenu.
type entry struct {
	label    string
	path     []string
	children []entry
}

var mainMenu = []entry{
	{label: msg.MenuSetup, path: []string{"setup"}},
	{label: msg.MenuQuick, path: []string{"tunnel"}},
	{label: msg.MenuRoutes, children: []entry{
		{label: msg.MenuList, path: []string{"list"}},
		{label: msg.MenuAdd, path: []string{"add"}},
		{label: msg.MenuEdit, path: []string{"edit"}},
		{label: msg.MenuRemove, path: []string{"rm"}},
		{label: msg.MenuCheck, path: []string{"check"}},
		{label: msg.MenuValidate, path: []string{"validate"}},
	}},
	{label: msg.MenuService, children: []entry{
		{label: msg.MenuServiceInstall, path: []string{"service", "install"}},
		{label: msg.MenuServiceStatus, path: []string{"service", "status"}},
		{label: msg.MenuServiceRestart, path: []string{"service", "restart"}},
		{label: msg.MenuServiceStart, path: []string{"service", "start"}},
		{label: msg.MenuServiceStop, path: []string{"service", "stop"}},
		{label: msg.MenuServiceEnable, path: []string{"service", "enable"}},
		{label: msg.MenuServiceDisable, path: []string{"service", "disable"}},
		{label: msg.MenuServiceLogs, path: []string{"service", "logs"}},
	}},
	{label: msg.MenuTunnel, children: []entry{
		{label: msg.MenuTunnelList, path: []string{"tunnels", "list"}},
		{label: msg.MenuTunnelCreate, path: []string{"tunnels", "create"}},
		{label: msg.MenuTunnelDelete, path: []string{"tunnels", "delete"}},
		{label: msg.MenuInit, path: []string{"init"}},
	}},
	{label: msg.MenuAccount, children: []entry{
		{label: msg.MenuInstall, path: []string{"install"}},
		{label: msg.MenuLogin, path: []string{"login"}},
		{label: msg.MenuTokenSet, path: []string{"token", "set"}},
		{label: msg.MenuTokenStatus, path: []string{"token", "status"}},
		{label: msg.MenuTokenClear, path: []string{"token", "clear"}},
	}},
	{label: msg.MenuUninstall, path: []string{"uninstall"}},
}

// runMenu shows the main menu until the user exits.
func runMenu(a *App, root *cobra.Command) error {
	return browse(a, root, msg.MenuTitle, mainMenu, msg.MenuExit)
}

// browse shows entries plus a leave option and dispatches the choice.
func browse(a *App, root *cobra.Command, title string, entries []entry, leave string) error {
	for {
		labels := make([]string, 0, len(entries)+1)
		for _, e := range entries {
			labels = append(labels, e.label)
		}
		i, err := a.Prompt.Select(title, append(labels, leave))
		if errors.Is(err, prompt.ErrAborted) || (err == nil && i == len(entries)) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := dispatch(a, root, entries[i]); err != nil {
			return err
		}
	}
}

// dispatch opens a submenu or runs a command, reporting command errors
// without leaving the menu.
func dispatch(a *App, root *cobra.Command, e entry) error {
	if e.children != nil {
		return browse(a, root, e.label, e.children, msg.MenuBack)
	}
	cmd, _, err := root.Find(e.path)
	if err != nil {
		return err
	}
	cmd.SetContext(root.Context())
	err = ensureRequirements(cmd, a)
	if err == nil {
		err = cmd.RunE(cmd, nil)
	}
	if err != nil {
		if errors.Is(err, prompt.ErrAborted) {
			return err
		}
		a.Printf(msg.ErrorFormat, err)
	}
	return nil
}
