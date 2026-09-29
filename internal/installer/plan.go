package installer

import (
	"fmt"

	"github.com/Jonathansl17/cloudfare/internal/sysexec"
)

const (
	releaseBaseURL = "https://github.com/cloudflare/cloudflared/releases/latest/download/"
	packageName    = "cloudflared"
	// BinaryPath is where the static binary is installed when no package fits.
	BinaryPath = "/usr/local/bin/cloudflared"
	binaryMode = "0755"
)

// plan describes how one family installs, detects and removes cloudflared.
type plan struct {
	assetFormat string
	archNames   map[string]string
	install     func(file string) []string
	query       []string
	remove      []string
}

var debArch = map[string]string{"amd64": "amd64", "arm64": "arm64", "arm": "arm", "386": "386"}
var rpmArch = map[string]string{"amd64": "x86_64", "arm64": "aarch64", "arm": "arm", "386": "386"}

var (
	debPlan = plan{
		assetFormat: "cloudflared-linux-%s.deb",
		archNames:   debArch,
		install:     func(f string) []string { return []string{"dpkg", "-i", f} },
		query:       []string{"dpkg", "-s", packageName},
		remove:      []string{"dpkg", "--purge", packageName},
	}
	rpmPlan = plan{
		assetFormat: "cloudflared-linux-%s.rpm",
		archNames:   rpmArch,
		install:     func(f string) []string { return []string{"rpm", "-Uvh", "--replacepkgs", f} },
		query:       []string{"rpm", "-q", packageName},
		remove:      []string{"rpm", "-e", packageName},
	}
	pacmanPlan = plan{
		install: func(string) []string { return []string{"pacman", "-S", "--noconfirm", "--needed", packageName} },
		query:   []string{"pacman", "-Q", packageName},
		remove:  []string{"pacman", "-Rns", "--noconfirm", packageName},
	}
	binaryPlan = plan{
		assetFormat: "cloudflared-linux-%s",
		archNames:   debArch,
		install:     func(f string) []string { return []string{"install", "-m", binaryMode, f, BinaryPath} },
		remove:      []string{"rm", "-f", BinaryPath},
	}
)

var plans = map[Family]plan{
	FamilyDebian:  debPlan,
	FamilyRHEL:    rpmPlan,
	FamilySUSE:    rpmPlan,
	FamilyArch:    pacmanPlan,
	FamilyAlpine:  binaryPlan,
	FamilyUnknown: binaryPlan,
}

func planFor(f Family) plan {
	if p, ok := plans[f]; ok {
		return p
	}
	return binaryPlan
}

// assetURL returns the release download for this plan, or "" when the plan
// installs from a distribution repository.
func (p plan) assetURL(goarch string) (string, error) {
	if p.assetFormat == "" {
		return "", nil
	}
	arch, ok := p.archNames[goarch]
	if !ok {
		return "", fmt.Errorf("unsupported CPU architecture %q", goarch)
	}
	return releaseBaseURL + fmt.Sprintf(p.assetFormat, arch), nil
}

func privileged(args []string) sysexec.Command {
	return sysexec.Command{Name: args[0], Args: args[1:], Privileged: true}
}
