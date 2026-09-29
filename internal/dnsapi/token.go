package dnsapi

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	tokenDir      = "cftun"
	tokenFile     = "token"
	tokenFilePerm = 0o600
	tokenDirPerm  = 0o700
)

// TokenStore keeps the API token in the user's config directory.
type TokenStore struct {
	Path string
	Env  func(string) string
}

// NewTokenStore locates the token file under the user config directory.
func NewTokenStore() (TokenStore, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return TokenStore{}, fmt.Errorf("locate config dir: %w", err)
	}
	return TokenStore{Path: filepath.Join(dir, tokenDir, tokenFile), Env: os.Getenv}, nil
}

// Load returns the token from the environment, then from the file, or "".
func (s TokenStore) Load() (string, error) {
	if token := s.Env(TokenEnv); token != "" {
		return token, nil
	}
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read token: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

// Save writes the token readable only by the current user.
func (s TokenStore) Save(token string) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), tokenDirPerm); err != nil {
		return fmt.Errorf("create token dir: %w", err)
	}
	if err := os.WriteFile(s.Path, []byte(token), tokenFilePerm); err != nil {
		return fmt.Errorf("write token: %w", err)
	}
	return nil
}

// Clear deletes the token file and its directory when empty.
func (s TokenStore) Clear() error {
	if err := os.Remove(s.Path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove token: %w", err)
	}
	err := os.Remove(filepath.Dir(s.Path))
	if err != nil && !errors.Is(err, fs.ErrNotExist) && !errors.Is(err, syscall.ENOTEMPTY) {
		return fmt.Errorf("remove token dir: %w", err)
	}
	return nil
}
