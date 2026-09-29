package routecmd

import (
	"github.com/spf13/cobra"

	"github.com/Jonathansl17/cftun/internal/cli/uikit"
	"github.com/Jonathansl17/cftun/internal/msg"
)

func (g *Group) validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   useValidate,
		Short: msg.ValidateShort,
		RunE:  g.validate,
	}
	return uikit.RequireCloudflared(cmd)
}

func (g *Group) validate(cmd *cobra.Command, _ []string) error {
	validation, err := g.deps.Validator.Validate(cmd.Context())
	if err != nil {
		return err
	}
	g.deps.Session.Printf(msg.InfoValidated, validation)
	return nil
}
