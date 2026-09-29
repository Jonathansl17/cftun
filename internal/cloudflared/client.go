// Package cloudflared wraps the cloudflared command line.
package cloudflared

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jonathansl17/cloudfare/internal/sysexec"
)

const (
	// Binary is the cloudflared executable name.
	Binary = "cloudflared"
	// HomeDir is the per-user directory holding cert.pem and tunnel credentials.
	HomeDir = ".cloudflared"
	// CertFile is written by `cloudflared login` inside HomeDir.
	CertFile        = "cert.pem"
	credentialsExt  = ".json"
	listOutputJSON  = "json"
	versionArg      = "--version"
	forceDeleteFlag = "-f"
)

// Tunnel is one entry of `cloudflared tunnel list`.
type Tunnel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Client runs cloudflared subcommands.
type Client struct {
	Runner sysexec.Runner
}

// Version returns the installed version, failing when cloudflared is absent.
func (c Client) Version(ctx context.Context) (string, error) {
	out, err := c.Runner.Output(ctx, c.cmd(versionArg))
	return strings.TrimSpace(out), err
}

// Login opens the browser authorization flow.
func (c Client) Login(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.cmd("tunnel", "login"))
}

// CreateTunnel creates a named tunnel and its credentials file.
func (c Client) CreateTunnel(ctx context.Context, name string) error {
	return c.Runner.Stream(ctx, c.cmd("tunnel", "create", name))
}

// DeleteTunnel removes stale connections and deletes the tunnel.
func (c Client) DeleteTunnel(ctx context.Context, tunnel string) error {
	if err := c.Runner.Stream(ctx, c.cmd("tunnel", "cleanup", tunnel)); err != nil {
		return err
	}
	return c.Runner.Stream(ctx, c.cmd("tunnel", "delete", forceDeleteFlag, tunnel))
}

// ListTunnels returns the account's active tunnels.
func (c Client) ListTunnels(ctx context.Context) ([]Tunnel, error) {
	out, err := c.Runner.Output(ctx, c.cmd("tunnel", "list", "--output", listOutputJSON))
	if err != nil {
		return nil, err
	}
	var tunnels []Tunnel
	if err := json.Unmarshal([]byte(out), &tunnels); err != nil {
		return nil, fmt.Errorf("parse tunnel list: %w", err)
	}
	return tunnels, nil
}

// RouteDNS points hostname at tunnel with a proxied CNAME.
func (c Client) RouteDNS(ctx context.Context, tunnel, hostname string) error {
	return c.Runner.Stream(ctx, c.cmd("tunnel", "route", "dns", tunnel, hostname))
}

// ValidateIngress checks the ingress rules of the config at path.
func (c Client) ValidateIngress(ctx context.Context, path string) (string, error) {
	return c.Runner.Output(ctx, c.cmd("tunnel", "--config", path, "ingress", "validate"))
}

// InstallService registers cloudflared as a system service.
func (c Client) InstallService(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.privileged("service", "install"))
}

// UninstallService removes the system service registration.
func (c Client) UninstallService(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.privileged("service", "uninstall"))
}

// CredentialsPath is where `tunnel create` writes the credentials of id.
func CredentialsPath(home, id string) string {
	return filepath.Join(home, HomeDir, id+credentialsExt)
}

// FindTunnel returns the tunnel whose name or ID equals ref.
func FindTunnel(tunnels []Tunnel, ref string) (Tunnel, bool) {
	for _, t := range tunnels {
		if t.Name == ref || t.ID == ref {
			return t, true
		}
	}
	return Tunnel{}, false
}

func (c Client) cmd(args ...string) sysexec.Command {
	return sysexec.Command{Name: Binary, Args: args}
}

func (c Client) privileged(args ...string) sysexec.Command {
	return sysexec.Command{Name: Binary, Args: args, Privileged: true}
}
