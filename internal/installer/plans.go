package installer

import "github.com/Jonathansl17/cftun/internal/paths"

var (
	debPlan = plan{
		assetFormat: debAssetFormat,
		install:     func(f string) []string { return []string{dpkgBinary, dpkgInstallFlag, f} },
		query:       []string{dpkgBinary, dpkgStatusFlag, packageName},
		remove:      []string{dpkgBinary, dpkgPurgeFlag, packageName},
	}
	rpmPlan = plan{
		assetFormat:   rpmAssetFormat,
		archOverrides: rpmArchOverrides,
		install:       func(f string) []string { return []string{rpmBinary, rpmUpgradeFlag, rpmReplaceFlag, f} },
		query:         []string{rpmBinary, rpmQueryFlag, packageName},
		remove:        []string{rpmBinary, rpmEraseFlag, packageName},
	}
	pacmanPlan = plan{
		install: func(string) []string {
			return []string{pacmanBinary, pacmanSyncFlag, pacmanNoConfirmFlag, pacmanNeededFlag, packageName}
		},
		query:  []string{pacmanBinary, pacmanQueryFlag, packageName},
		remove: []string{pacmanBinary, pacmanRemoveFlags, pacmanNoConfirmFlag, packageName},
	}
	binaryPlan = plan{
		assetFormat: binaryAssetFormat,
		install: func(f string) []string {
			return []string{installBinary, installModeFlag, binaryMode, f, paths.InstallPath}
		},
		remove: []string{removeBinary, removeForceFlag, paths.InstallPath},
	}
)

var plans = map[Family]plan{
	FamilyDebian:  debPlan,
	FamilyRHEL:    rpmPlan,
	FamilySUSE:    rpmPlan,
	FamilyArch:    pacmanPlan,
	FamilyAlpine:  binaryPlan,
	FamilyUnknown: binaryPlan,
}
