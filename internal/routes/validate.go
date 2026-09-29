package routes

import "context"

func (v Validator) Validate(ctx context.Context) (string, error) {
	return v.Tunnels.ValidateIngress(ctx, v.Path.Path())
}
