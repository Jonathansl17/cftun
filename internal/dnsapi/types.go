package dnsapi

import (
	"encoding/json"
	"net/http"
)

type TokenSource interface {
	Load() (string, error)
}

type Client struct {
	HTTP    *http.Client
	BaseURL string
	Tokens  TokenSource
}

type TokenStore struct {
	Path string
	Env  func(string) string
}

type APIError struct {
	Op  string
	Err error
}

type envelope struct {
	Success bool              `json:"success"`
	Errors  []json.RawMessage `json:"errors"`
	Result  []struct {
		ID string `json:"id"`
	} `json:"result"`
}
