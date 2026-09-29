package msg

const (
	HintConnectionRefused = "nothing listens on that port; start your app"
	HintDNS1016           = "error 1016: DNS does not point to the tunnel; remove and add the route again"
	HintTunnelDown        = "the tunnel cannot reach the app; check the port and \"cftun service status\""
	HintGeneric           = "unexpected answer; check \"cftun service logs\""
)

const (
	StepStopService      = "stop service"
	StepDeleteDNS        = "delete DNS record"
	StepDeleteTunnel     = "delete tunnel"
	StepUninstallService = "unregister service"
	StepRemovePackage    = "remove package"
	StepRemoveFiles      = "remove files"
	StepRemoveTemp       = "remove temporary tunnel configs"
)
