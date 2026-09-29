package dnsapi

const (
	DefaultBaseURL = "https://api.cloudflare.com/client/v4"
	TokenEnv       = "CLOUDFLARE_API_TOKEN"

	recordType     = "CNAME"
	minZoneLabels  = 2
	labelSeparator = "."

	zonesPathFormat     = "/zones?%s"
	recordsPathFormat   = "/zones/%s/dns_records?%s"
	recordPathFormat    = "/zones/%s/dns_records/%s"
	queryType           = "type"
	queryName           = "name"
	headerAuthorization = "Authorization"
	bearerFormat        = "Bearer %s"

	opBuildRequest = "build request"
	opSend         = "send"
	opStatus       = "status"
	opDecode       = "decode"
	opDrain        = "drain response"
	opResponse     = "response"

	apiErrorFormat    = "cloudflare api %s: %v"
	statusErrorFormat = "unexpected status %s"
	failureFormat     = "%s: %s"
	zoneErrorFormat   = "%s: %w"

	tokenDir      = "cftun"
	tokenFile     = "token"
	tokenFilePerm = 0o600
	tokenDirPerm  = 0o700
)
