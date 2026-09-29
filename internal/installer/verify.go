package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func VerifySHA256(path, digest string) error {
	if digest == "" {
		return fmt.Errorf(missingDigestFormat, ErrMissingDigest, path)
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf(openFileFormat, path, err)
	}
	sum, err := hashReader(path, f)
	if closeErr := f.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf(hashFileFormat, path, closeErr)
	}
	if err != nil {
		return err
	}
	if !strings.EqualFold(sum, digest) {
		return fmt.Errorf(checksumFormat, ErrChecksumMismatch, digest, sum)
	}
	return nil
}

func hashReader(path string, r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf(hashFileFormat, path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
