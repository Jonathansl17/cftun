package installer

import "github.com/Jonathansl17/cftun/internal/paths"

const (
	FamilyUnknown Family = iota
	FamilyDebian
	FamilyRHEL
	FamilySUSE
	FamilyArch
	FamilyAlpine
)

const (
	latestReleaseURL = "https://api.github.com/repos/cloudflare/cloudflared/releases/latest"
	acceptHeader     = "Accept"
	acceptGitHubJSON = "application/vnd.github+json"
	apiVersionHeader = "X-GitHub-Api-Version"
	apiVersion       = "2022-11-28"
	digestPrefix     = "sha256:"
	httpsScheme      = "https"
	maxReleaseBytes  = 4 << 20
	maxRedirects     = 10

	packageName  = paths.CloudflaredBinary
	binaryMode   = "0755"
	downloadPerm = 0o644

	downloadDirPattern = "cftun-download-*"

	queryNotInstalledExitCode = 1

	unsupportedArchFormat  = "%w %q"
	assetNotFoundFormat    = "%w %q"
	missingDigestFormat    = "%w %q"
	insecureURLFormat      = "%w %q"
	insecureRedirectFormat = "%w %q"
	checksumFormat         = "%w: want %s, got %s"
	releaseRequestFormat   = "build release request: %w"
	releaseFetchFormat     = "fetch latest release: %w"
	releaseStatusFormat    = "fetch latest release: unexpected status %s"
	releaseDecodeFormat    = "decode latest release: %w"
	buildRequestFormat     = "build request: %w"
	downloadFormat         = "download %s: %w"
	downloadStatusFormat   = "download %s: unexpected status %s"
	createFileFormat       = "create %s: %w"
	saveFileFormat         = "save %s: %w"
	openFileFormat         = "open %s: %w"
	hashFileFormat         = "hash %s: %w"
	parseURLFormat         = "parse url %q: %w"
	queryFailedFormat      = "query package %s: %w"
	readOSReleaseFormat    = "read %s: %w"

	osReleaseIDKey     = "ID"
	osReleaseIDLikeKey = "ID_LIKE"
	osReleaseQuotes    = `"'`

	debAssetFormat    = "cloudflared-linux-%s.deb"
	rpmAssetFormat    = "cloudflared-linux-%s.rpm"
	binaryAssetFormat = "cloudflared-linux-%s"

	dpkgBinary          = "dpkg"
	dpkgInstallFlag     = "-i"
	dpkgStatusFlag      = "-s"
	dpkgPurgeFlag       = "--purge"
	rpmBinary           = "rpm"
	rpmUpgradeFlag      = "-Uvh"
	rpmReplaceFlag      = "--replacepkgs"
	rpmQueryFlag        = "-q"
	rpmEraseFlag        = "-e"
	pacmanBinary        = "pacman"
	pacmanSyncFlag      = "-S"
	pacmanNoConfirmFlag = "--noconfirm"
	pacmanNeededFlag    = "--needed"
	pacmanQueryFlag     = "-Q"
	pacmanRemoveFlags   = "-Rns"
	installBinary       = "install"
	installModeFlag     = "-m"
	removeBinary        = "rm"
	removeForceFlag     = "-f"
)

var supportedArches = []string{"amd64", "arm64", "arm", "386"}

var rpmArchOverrides = map[string]string{"amd64": "x86_64", "arm64": "aarch64"}

var familyByID = map[string]Family{
	"debian":   FamilyDebian,
	"ubuntu":   FamilyDebian,
	"rhel":     FamilyRHEL,
	"fedora":   FamilyRHEL,
	"centos":   FamilyRHEL,
	"suse":     FamilySUSE,
	"opensuse": FamilySUSE,
	"arch":     FamilyArch,
	"alpine":   FamilyAlpine,
}

var familyNames = map[Family]string{
	FamilyUnknown: "unknown",
	FamilyDebian:  "debian",
	FamilyRHEL:    "rhel",
	FamilySUSE:    "suse",
	FamilyArch:    "arch",
	FamilyAlpine:  "alpine",
}
