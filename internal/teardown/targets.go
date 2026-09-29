package teardown

func targetsFor(params Params) []Target {
	targets := append([]Target{}, systemTargets...)
	return append(targets,
		Target{Path: params.ConfigPath, Kind: TargetFile},
		Target{Path: params.BackupPath, Kind: TargetFile},
		Target{Path: params.UserDir, Kind: TargetDirectory},
	)
}
