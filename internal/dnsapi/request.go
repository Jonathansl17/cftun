package dnsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c Client) ids(ctx context.Context, method, path string) ([]string, error) {
	env, err := c.call(ctx, method, path)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(env.Result))
	for _, r := range env.Result {
		ids = append(ids, r.ID)
	}
	return ids, nil
}

func (c Client) call(ctx context.Context, method, path string) (*envelope, error) {
	token, err := c.token()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, nil)
	if err != nil {
		return nil, apiError(opBuildRequest, err)
	}
	req.Header.Set(headerAuthorization, fmt.Sprintf(bearerFormat, token))
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, apiError(opSend, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, apiError(opStatus, fmt.Errorf(statusErrorFormat, resp.Status))
	}
	if method == http.MethodDelete {
		return drain(resp.Body)
	}
	return decode(resp)
}

func (c Client) token() (string, error) {
	token, err := c.Tokens.Load()
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", ErrNoToken
	}
	return token, nil
}

func decode(resp *http.Response) (*envelope, error) {
	var env envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, apiError(opDecode, err)
	}
	if !env.Success {
		return nil, apiError(opResponse, fmt.Errorf(failureFormat, resp.Status, env.Errors))
	}
	return &env, nil
}

func drain(body io.Reader) (*envelope, error) {
	if _, err := io.Copy(io.Discard, body); err != nil {
		return nil, apiError(opDrain, err)
	}
	return &envelope{}, nil
}

func apiError(op string, err error) error {
	return &APIError{Op: op, Err: err}
}
