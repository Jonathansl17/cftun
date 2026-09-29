package routes

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/apperr"
	"github.com/Jonathansl17/cftun/internal/cloudflared"
	"github.com/Jonathansl17/cftun/internal/ingress"
)

func (i Initializer) Init(ctx context.Context, tunnelRef string) (string, error) {
	home, err := i.Home.Home()
	if err != nil {
		return "", err
	}
	tunnels, err := i.Catalog.ListTunnels(ctx)
	if err != nil {
		return "", err
	}
	tunnel, ok := cloudflared.FindTunnel(tunnels, tunnelRef)
	if !ok {
		return "", apperr.Wrap(tunnelRef, ErrTunnelNotFound)
	}
	creds := cloudflared.CredentialsPath(home, tunnel.ID)
	if !i.Files.Exists(creds) {
		return "", apperr.Wrap(creds, ErrCredentialsMissing)
	}
	a := applier{store: i.Store, checker: i.Checker}
	rollback := a.discard
	if i.Store.Exists() {
		rollback = a.restore
	}
	return a.commit(ctx, ingress.New(tunnel.ID, creds), rollback)
}
