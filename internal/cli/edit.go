package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/msg"
	"github.com/Jonathansl17/cftun/internal/routes"
)

func newEditCmd(a *App) *cobra.Command {
	var host, newHost, port string
	var opts routes.Options
	cmd := &cobra.Command{
		Use:         "edit",
		Short:       msg.EditShort,
		Annotations: needsCloudflared,
		RunE: func(cmd *cobra.Command, _ []string) error {
			rule, err := pickRule(a, host)
			if err != nil {
				return err
			}
			updated, err := askEdit(a, rule, newHost, port)
			if err != nil {
				return err
			}
			out, err := a.Editor.Edit(cmd.Context(), rule.Hostname, updated, opts)
			printOutcome(a, rule.Hostname, out)
			if err != nil {
				return present(err)
			}
			a.Printf(msg.InfoAdded, updated.Hostname, updated.Service)
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "", msg.FlagHost)
	cmd.Flags().StringVar(&newHost, "new-host", "", msg.FlagNewHost)
	cmd.Flags().StringVar(&port, "port", "", msg.FlagPort)
	bindRouteOptions(cmd, &opts)
	return cmd
}

func askEdit(a *App, rule ingress.Rule, newHost, port string) (ingress.Rule, error) {
	if newHost == "" && port == "" {
		var err error
		if newHost, port, err = askEditInteractive(a, rule); err != nil {
			return ingress.Rule{}, err
		}
	}
	hostname, err := askHostname(a, orDefault(newHost, rule.Hostname), "")
	if err != nil {
		return ingress.Rule{}, err
	}
	service := rule.Service
	if port != "" {
		p, err := askPort(a, port, "")
		if err != nil {
			return ingress.Rule{}, err
		}
		service = ingress.LocalService(p)
	}
	return ingress.Rule{Hostname: hostname, Service: service}, nil
}

func askEditInteractive(a *App, rule ingress.Rule) (string, string, error) {
	currentPort := ""
	if p, ok := ingress.Port(rule.Service); ok {
		currentPort = strconv.Itoa(p)
	}
	host, err := keepIfBlank(a, msg.PromptNewHostname, rule.Hostname)
	if err != nil {
		return "", "", err
	}
	port, err := keepIfBlank(a, msg.PromptNewPort, currentPort)
	return host, port, err
}

func keepIfBlank(a *App, label, current string) (string, error) {
	answer, err := a.Prompt.Ask(label+msg.CurrentValue(current), func(string) error { return nil })
	if err != nil {
		return "", err
	}
	return orDefault(answer, current), nil
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
