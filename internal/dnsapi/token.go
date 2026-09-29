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

func NewTokenStore() (TokenStore, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return TokenStore{}, fmt.Errorf(locateConfigDirFormat, err)
	}
	return TokenStore{Path: filepath.Join(dir, tokenDir, tokenFile), Env: os.Getenv}, nil
}

func (s TokenStore) Location() string {
	return s.Path
}

func (s TokenStore) Load() (string, error) {
	if token := s.Env(TokenEnv); token != "" {
		return token, nil
	}
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf(readTokenFormat, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func (s TokenStore) Save(token string) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), tokenDirPerm); err != nil {
		return fmt.Errorf(createTokenDirFormat, err)
	}
	if err := os.WriteFile(s.Path, []byte(token), tokenFilePerm); err != nil {
		return fmt.Errorf(writeTokenFormat, err)
	}
	if err := os.Chmod(s.Path, tokenFilePerm); err != nil {
		return fmt.Errorf(restrictTokenFormat, err)
	}
	return nil
}

func (s TokenStore) Clear() error {
	if err := os.Remove(s.Path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf(removeTokenFormat, err)
	}
	err := os.Remove(filepath.Dir(s.Path))
	if err != nil && !errors.Is(err, fs.ErrNotExist) && !errors.Is(err, syscall.ENOTEMPTY) {
		return fmt.Errorf(removeTokenDirFormat, err)
	}
	return nil
}
