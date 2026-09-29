package routes

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Jonathansl17/cloudfare/internal/cloudflared"
	"github.com/Jonathansl17/cloudfare/internal/ingress"
)

var (
	// ErrConfigExists protects an existing config from being overwritten.
	ErrConfigExists = errors.New("config already exists, pass --force to overwrite it")
	// ErrTunnelNotFound reports a tunnel name or ID that the account lacks.
	ErrTunnelNotFound = errors.New("tunnel not found, create it with `cftun tunnel create`")
)

// Init writes a fresh config pointing at tunnelRef, whose credentials live
// under home.
func (m Manager) Init(ctx context.Context, tunnelRef, home string, force bool) error {
	if m.Store.Exists() && !force {
		return ErrConfigExists
	}
	tunnels, err := m.Tunnels.ListTunnels(ctx)
	if err != nil {
		return err
	}
	tunnel, ok := cloudflared.FindTunnel(tunnels, tunnelRef)
	if !ok {
		return fmt.Errorf("%s: %w", tunnelRef, ErrTunnelNotFound)
	}
	creds := cloudflared.CredentialsPath(home, tunnel.ID)
	if _, err := os.Stat(creds); err != nil {
		return fmt.Errorf("credentials file: %w", err)
	}
	if err := m.Store.Save(ctx, ingress.New(tunnel.ID, creds)); err != nil {
		return err
	}
	return m.Validate(ctx)
}
