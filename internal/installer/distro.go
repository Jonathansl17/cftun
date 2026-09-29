// Package installer installs and removes cloudflared on any Linux distribution.
package installer

import (
	"bufio"
	"io"
	"strings"
)

// OSReleasePath is the standard distribution identification file.
const OSReleasePath = "/etc/os-release"

// Family groups distributions that share a package format.
type Family int

// Supported families. Unknown falls back to the static binary.
const (
	FamilyUnknown Family = iota
	FamilyDebian
	FamilyRHEL
	FamilySUSE
	FamilyArch
	FamilyAlpine
)

var familyByID = map[string]Family{
	"debian":   FamilyDebian,
	"ubuntu":   FamilyDebian,
	"rhel":     FamilyRHEL,
	"fedora":   FamilyRHEL,
	"centos":   FamilyRHEL,
	"suse":     FamilySUSE,
	"opensuse": FamilySUSE,
	"arch":     FamilyArch,
	"alpine":   FamilyAlpine,
}

var familyNames = map[Family]string{
	FamilyUnknown: "unknown",
	FamilyDebian:  "debian",
	FamilyRHEL:    "rhel",
	FamilySUSE:    "suse",
	FamilyArch:    "arch",
	FamilyAlpine:  "alpine",
}

func (f Family) String() string {
	return familyNames[f]
}

// DetectFamily reads an os-release file and matches ID, then ID_LIKE.
func DetectFamily(r io.Reader) (Family, error) {
	fields, err := parseOSRelease(r)
	if err != nil {
		return FamilyUnknown, err
	}
	candidates := append([]string{fields["ID"]}, strings.Fields(fields["ID_LIKE"])...)
	for _, id := range candidates {
		if f, ok := familyByID[id]; ok {
			return f, nil
		}
	}
	return FamilyUnknown, nil
}

func parseOSRelease(r io.Reader) (map[string]string, error) {
	fields := map[string]string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		key, value, ok := strings.Cut(strings.TrimSpace(scanner.Text()), "=")
		if ok {
			fields[key] = strings.ToLower(strings.Trim(value, `"'`))
		}
	}
	return fields, scanner.Err()
}
