package cloudflared

import "github.com/Jonathansl17/cftun/internal/sysexec"

type Tunnel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	Runner sysexec.Runner
}
