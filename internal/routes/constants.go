package routes

const (
	restoreBackupFormat   = "restore backup: %w"
	removeNewConfigFormat = "remove new config: %w"
	rollbackFormat        = "config rejected, previous version restored: %v"
	dnsRollbackFormat     = "dns route failed, previous config restored: %v"
)
