package installer

import (
	"bufio"
	"io"
	"strings"
)

func (f Family) String() string {
	return familyNames[f]
}

func DetectFamily(r io.Reader) (Family, error) {
	fields, err := parseOSRelease(r)
	if err != nil {
		return FamilyUnknown, err
	}
	candidates := append([]string{fields[osReleaseIDKey]}, strings.Fields(fields[osReleaseIDLikeKey])...)
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
			fields[key] = strings.ToLower(strings.Trim(value, osReleaseQuotes))
		}
	}
	return fields, scanner.Err()
}
