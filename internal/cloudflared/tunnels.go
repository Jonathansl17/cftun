package cloudflared

import (
	"path/filepath"

	"github.com/Jonathansl17/cftun/internal/paths"
)

func CredentialsPath(home, id string) string {
	return filepath.Join(home, paths.UserDirName, id+credentialsExt)
}

func FindTunnel(tunnels []Tunnel, ref string) (Tunnel, bool) {
	for _, t := range tunnels {
		if t.Name == ref || t.ID == ref {
			return t, true
		}
	}
	return Tunnel{}, false
}
