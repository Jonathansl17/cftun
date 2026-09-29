package ingress

import "regexp"

const (
	keyTunnel      = "tunnel"
	keyCredentials = "credentials-file"
	keyIngress     = "ingress"
	keyHostname    = "hostname"
	keyService     = "service"

	yamlIndent = 2
	pairStride = 2

	CatchAllService = "http_status:404"

	MinPort            = 1
	MaxPort            = 65535
	localServiceFormat = "http://localhost:%d"
)

var hostnamePattern = regexp.MustCompile(`^([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$`)
