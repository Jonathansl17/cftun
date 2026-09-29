package store

const (
	backupSuffix    = ".bak"
	filePerm        = 0o644
	fileModeFormat  = "%04o"
	tempFilePattern = "cftun-*"

	removeTempFormat = "remove temp file: %w"

	installBinary = "install"
	installDirs   = "-D"
	installMode   = "-m"

	removeBinary = "rm"
	removeForce  = "-f"
	endOfOptions = "--"
)
