# cftun

A CLI to run a single Cloudflare Tunnel on a Linux server and publish local
apps on your domain, without editing YAML or remembering `cloudflared` flags.

One tunnel carries every app. Each route maps a public hostname to a local
port:

```
https://api.example.com  ->  http://localhost:5000
https://example.com      ->  http://localhost:5174
```

Cloudflare terminates HTTPS at its edge, so your apps only need to listen on
plain HTTP on `localhost`. No Nginx, no open ports, no certificates.

## Install

### One-liner (recommended, pinned to the latest stable tag)

```sh
curl -fsSL https://raw.githubusercontent.com/Jonathansl17/cftun/v1.0.1/install.sh | CFTUN_VERSION=v1.0.1 bash
```

Tags are immutable, so this URL is not affected by the GitHub raw CDN cache.

### Track the latest release

```sh
curl -fsSL https://raw.githubusercontent.com/Jonathansl17/cftun/master/install.sh | bash
```

The installer detects the CPU architecture (amd64, arm64, arm, 386),
downloads the matching binary from the
[releases](https://github.com/Jonathansl17/cftun/releases), verifies its
SHA-256 checksum and installs it to `/usr/local/bin/cftun` (using sudo when
needed). It needs `curl` or `wget` and nothing else: no Go, no runtime.

| Variable | Default | Purpose |
|---|---|---|
| `CFTUN_VERSION` | `latest` | Release tag to install |
| `CFTUN_INSTALL_DIR` | `/usr/local/bin` | Target directory |

Re-run the installer to update. See [Uninstall](#uninstall) to remove it.

### Build from source

Requirements: Go 1.22 or newer and git.

```sh
git clone https://github.com/Jonathansl17/cftun.git
cd cftun
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(git describe --tags --always)" -o cftun ./cmd/cftun
sudo install -m 0755 cftun /usr/local/bin/cftun
cftun --version
```

`CGO_ENABLED=0` produces a static binary that runs on any distribution,
including musl-based ones like Alpine.

To build on your machine for a server with another CPU, set `GOARCH`
(`amd64`, `arm64`, `arm` or `386`) and copy the binary over:

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -o cftun ./cmd/cftun
scp cftun user@server:~
ssh user@server 'sudo install -m 0755 ~/cftun /usr/local/bin/cftun'
```

To build every architecture at once into `dist/` with checksums:

```sh
scripts/build-release.sh v1.0.1
```

## Quick start

On a fresh machine, run the guided setup as your normal user (not with sudo;
cftun asks for sudo only when it needs it):

```sh
cftun setup
```

It walks through every step and skips what is already done:

1. Installs `cloudflared` with the native method of your distribution.
2. Runs `cloudflared tunnel login` so you can pick your domain in the browser.
3. Creates the tunnel (or reuses it) and writes `/etc/cloudflared/config.yml`.
4. Registers the system service, enables it on boot and starts it.
5. Optionally saves an API token (see [API token](#api-token)).
6. Lets you add your first routes.

After that, run `cftun` without arguments for the interactive menu, or use
the subcommands below.

## Commands

| Command | What it does |
|---|---|
| `cftun` | Interactive menu |
| `cftun setup` | Guided setup from zero to a working tunnel |
| `cftun add [--host H] [--port P]` | Add a route, create its DNS record, restart the service |
| `cftun rm [--host H]` | Remove a route and its DNS record, restart the service |
| `cftun edit [--host H] [--new-host N] [--port P]` | Change a route's hostname or port |
| `cftun list` | List routes and whether their local port answers |
| `cftun check` | Test every route locally and from the internet, with hints |
| `cftun validate` | Validate the ingress rules with cloudflared |
| `cftun install [--force]` | Install cloudflared |
| `cftun login` | Authorize cloudflared with your account |
| `cftun tunnel create\|list\|delete` | Manage tunnels |
| `cftun init [--tunnel T] [--force]` | Write the config for an existing tunnel |
| `cftun service install` | Register, enable and start the service |
| `cftun service start\|stop\|restart\|status\|enable\|disable\|logs` | Control the service |
| `cftun token set\|status\|clear` | Manage the API token |
| `cftun uninstall [-y]` | Remove everything, leaving the machine as it was |

Missing values are asked interactively, so `cftun add` alone asks for the
hostname and the port. `add`, `rm` and `edit` accept `--no-dns` and
`--no-restart` to skip those steps.

Every change backs up the config to `config.yml.bak` and runs
`cloudflared tunnel ingress validate`. If validation fails, the backup is
restored. The catch-all `http_status:404` rule always stays last.

## API token

`cloudflared` can create DNS records but cannot delete them. To let `rm`,
`edit` and `uninstall` remove old CNAME records:

1. Open <https://dash.cloudflare.com/profile/api-tokens>.
2. Create a token from the **Edit zone DNS** template for your zone.
3. Run `cftun token set` and paste it.

The token is stored in `~/.config/cftun/token` with mode `0600`. The
`CLOUDFLARE_API_TOKEN` environment variable takes precedence. Without a
token, cftun tells you which record to delete in the dashboard.

## Supported systems

| Distribution family | Install method |
|---|---|
| Debian, Ubuntu and derivatives | `.deb` from the cloudflared releases |
| Fedora, RHEL, CentOS, Rocky, Alma | `.rpm` from the cloudflared releases |
| openSUSE, SLES | `.rpm` from the cloudflared releases |
| Arch, Manjaro, EndeavourOS | `pacman` |
| Alpine and anything else | static binary in `/usr/local/bin` |

CPU architectures: amd64, arm64, arm and 386. Services work with systemd,
OpenRC and SysV init.

## Configuration

| Setting | Default | Override |
|---|---|---|
| cloudflared config | `/etc/cloudflared/config.yml` | `--config` or `CFTUN_CONFIG` |
| API token | `~/.config/cftun/token` | `CLOUDFLARE_API_TOKEN` |

## Uninstall

One-liner that removes cloudflared, everything cftun created and the cftun
binary itself:

```sh
curl -fsSL https://raw.githubusercontent.com/Jonathansl17/cftun/master/uninstall.sh | bash
```

| Variable | Default | Purpose |
|---|---|---|
| `CFTUN_YES` | `0` | Set to `1` to skip the confirmation |
| `CFTUN_KEEP_CLOUDFLARED` | `0` | Set to `1` to remove only the cftun binary |
| `CFTUN_INSTALL_DIR` | `/usr/local/bin` | Where cftun was installed |

To keep cftun and only remove cloudflared, run `cftun uninstall` instead.
Either way, it stops and unregisters the service, deletes the DNS record of
every route (with an API token), deletes the tunnel, removes the package and
deletes `/etc/cloudflared`, `~/.cloudflared`, the logs and the saved token. If
a step fails, it keeps going and reports every failure at the end; the cftun
binary is only removed once the cleanup succeeded, so you can re-run it.

## Troubleshooting

`cftun check` detects these and prints a hint:

- **Connection refused**: nothing listens on the local port; start your app.
- **Error 1016**: the DNS record does not point to the tunnel; remove and add
  the route again.
- **502 / 530**: the tunnel cannot reach the app; check the port and
  `cftun service status`.
