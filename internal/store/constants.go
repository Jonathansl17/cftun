package store

const (
	backupSuffix    = ".bak"
	filePerm        = 0o644
	fileModeFormat  = "%04o"
	tempFilePattern = "cftun-*"

	readFileFormat   = "read %s: %w"
	readBackupFormat = "read backup: %w"
	writeFileFormat  = "write %s: %w"
	removeFileFormat = "remove %s: %w"
	createTempFormat = "create temp file: %w"
	writeTempFormat  = "write temp file: %w"
	closeTempFormat  = "close temp file: %w"
	removeTempFormat = "remove temp file: %w"

	installBinary = "install"
	installDirs   = "-D"
	installMode   = "-m"

	removeBinary = "rm"
	removeForce  = "-f"
	endOfOptions = "--"
)
