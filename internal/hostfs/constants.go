package hostfs

import "os"

const (
	tempFilePerm os.FileMode = 0o600

	createTempFormat        = "create temp file: %w"
	closeTempFormat         = "close temp file: %w"
	removeTempFormat        = "remove temp file: %w"
	globTempFormat          = "glob temp files: %w"
	writeTempFormat         = "write temp file: %w"
	inspectFormat           = "inspect %s: %w"
	locateExecutableFormat  = "locate running executable: %w"
	resolveExecutableFormat = "resolve executable %s: %w"
	locateHomeFormat        = "locate home: %w"
)
