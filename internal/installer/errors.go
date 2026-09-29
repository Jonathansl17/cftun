package installer

import "errors"

var (
	ErrUnsupportedArch  = errors.New("unsupported CPU architecture")
	ErrAssetNotFound    = errors.New("release asset not found")
	ErrMissingDigest    = errors.New("release asset has no sha256 digest")
	ErrInsecureURL      = errors.New("download url is not https")
	ErrInsecureRedirect = errors.New("download redirected to a non-https url")
	ErrChecksumMismatch = errors.New("sha256 checksum mismatch")
)
