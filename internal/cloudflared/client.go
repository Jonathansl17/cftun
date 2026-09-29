package cloudflared

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jonathansl17/cftun/internal/sysexec"
)

const (
	Binary           = "cloudflared"
	HomeDir          = ".cloudflared"
	CertFile         = "cert.pem"
	credentialsExt   = ".json"
	listOutputJSON   = "json"
	versionArg       = "--version"
	forceDeleteFlag  = "-f"
	quickGracePeriod = "1s"
)

type Tunnel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	Runner sysexec.Runner
}

func (c Client) Version(ctx context.Context) (string, error) {
	out, err := c.Runner.Output(ctx, c.cmd(versionArg))
	return strings.TrimSpace(out), err
}

func (c Client) Login(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.cmd("tunnel", "login"))
}

func (c Client) CreateTunnel(ctx context.Context, name string) error {
	return c.Runner.Stream(ctx, c.cmd("tunnel", "create", name))
}

func (c Client) DeleteTunnel(ctx context.Context, tunnel string) error {
	if err := c.Runner.Stream(ctx, c.cmd("tunnel", "cleanup", tunnel)); err != nil {
		return err
	}
	return c.Runner.Stream(ctx, c.cmd("tunnel", "delete", forceDeleteFlag, tunnel))
}

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

func (c Client) RouteDNS(ctx context.Context, tunnel, hostname string) error {
	return c.Runner.Stream(ctx, c.cmd("tunnel", "route", "dns", tunnel, hostname))
}

func (c Client) ValidateIngress(ctx context.Context, path string) (string, error) {
	return c.Runner.Output(ctx, c.cmd("tunnel", "--config", path, "ingress", "validate"))
}

func (c Client) QuickTunnel(ctx context.Context, url, emptyConfig string) error {
	return c.Runner.Stream(ctx, c.cmd("tunnel", "--config", emptyConfig, "--no-autoupdate",
		"--grace-period", quickGracePeriod, "--url", url))
}

func (c Client) InstallService(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.privileged("service", "install"))
}

func (c Client) UninstallService(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.privileged("service", "uninstall"))
}

func CredentialsPath(home, id string) string {
	return filepath.Join(home, HomeDir, id+credentialsExt)
}

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
