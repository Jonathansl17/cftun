package cloudflared

import "context"

func (c Client) QuickTunnel(ctx context.Context, url, emptyConfig string) error {
	return c.Runner.Stream(ctx, c.cmd(subTunnel, flagConfig, emptyConfig, flagNoAutoupdate,
		flagGracePeriod, quickGracePeriod, flagURL, url))
}
