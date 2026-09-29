package cloudflared

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jonathansl17/cftun/internal/paths"
	"github.com/Jonathansl17/cftun/internal/sysexec"
)

func (c Client) Version(ctx context.Context) (string, error) {
	out, err := c.Runner.Output(ctx, c.cmd(flagVersion))
	return strings.TrimSpace(out), err
}

func (c Client) Login(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.cmd(subTunnel, subLogin))
}

func (c Client) CreateTunnel(ctx context.Context, name string) error {
	return c.Runner.Stream(ctx, c.cmd(subTunnel, subCreate, name))
}

func (c Client) DeleteTunnel(ctx context.Context, tunnel string) error {
	if err := c.Runner.Stream(ctx, c.cmd(subTunnel, subCleanup, tunnel)); err != nil {
		return err
	}
	return c.Runner.Stream(ctx, c.cmd(subTunnel, subDelete, flagForce, tunnel))
}

func (c Client) ListTunnels(ctx context.Context) ([]Tunnel, error) {
	out, err := c.Runner.Output(ctx, c.cmd(subTunnel, subList, flagOutput, outputJSON))
	if err != nil {
		return nil, err
	}
	var tunnels []Tunnel
	if err := json.Unmarshal([]byte(out), &tunnels); err != nil {
		return nil, fmt.Errorf(parseTunnelListFormat, err)
	}
	return tunnels, nil
}

func (c Client) RouteDNS(ctx context.Context, tunnel, hostname string) error {
	return c.Runner.Stream(ctx, c.cmd(subTunnel, subRoute, subDNS, tunnel, hostname))
}

func (c Client) ValidateIngress(ctx context.Context, path string) (string, error) {
	return c.Runner.Output(ctx, c.cmd(subTunnel, flagConfig, path, subIngress, subValidate))
}

func (c Client) InstallService(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.privileged(subService, subInstall))
}

func (c Client) UninstallService(ctx context.Context) error {
	return c.Runner.Stream(ctx, c.privileged(subService, subUninstall))
}

func (c Client) cmd(args ...string) sysexec.Command {
	return sysexec.Command{Name: paths.CloudflaredBinary, Args: args}
}

func (c Client) privileged(args ...string) sysexec.Command {
	return sysexec.Command{Name: paths.CloudflaredBinary, Args: args, Privileged: true}
}
