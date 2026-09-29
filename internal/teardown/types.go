package teardown

import (
	"context"

	"github.com/Jonathansl17/cftun/internal/hostfs"
	"github.com/Jonathansl17/cftun/internal/ingress"
	"github.com/Jonathansl17/cftun/internal/routes"
	"github.com/Jonathansl17/cftun/internal/sysexec"
)

type Params struct {
	ConfigPath string
	BackupPath string
	UserDir    string
}

type TargetKind int

const (
	TargetDirectory TargetKind = iota
	TargetFile
)

type Target struct {
	Path string
	Kind TargetKind
}

type Step int

const (
	StepStopService Step = iota
	StepDeleteDNS
	StepDeleteTunnel
	StepUninstallService
	StepRemovePackage
	StepRemoveFiles
	StepRemoveTemp
)

type EventKind int

const (
	EventStepFailed EventKind = iota
	EventNoConfig
	EventRemoving
	EventDNSManual
	EventDNSDeleted
)

type Event struct {
	Kind    EventKind
	Step    Step
	Subject string
	Count   int
	Err     error
}

type Observer interface {
	Observe(Event)
}

type ConfigReader interface {
	Load() (*ingress.Document, error)
}

type TunnelOps interface {
	Version(ctx context.Context) (string, error)
	DeleteTunnel(ctx context.Context, tunnel string) error
	UninstallService(ctx context.Context) error
}

type ServiceControl interface {
	Stop(ctx context.Context) error
	Registered() bool
}

type PackageRemover interface {
	Uninstall(ctx context.Context) error
}

type DNSCleaning interface {
	Delete(ctx context.Context, host string) (routes.DNSRemoval, error)
}

type FileRemover interface {
	RemoveAll(ctx context.Context, targets []Target) error
}

type TempLocator interface {
	TempMatches(pattern string) ([]string, error)
}

type KindInspector interface {
	KindOf(path string) (hostfs.Kind, error)
}

type Procedure struct {
	Config   ConfigReader
	Tunnels  TunnelOps
	Service  ServiceControl
	Package  PackageRemover
	DNS      DNSCleaning
	Files    FileRemover
	Temp     TempLocator
	Observer Observer
}

type SafeRemover struct {
	Runner    sysexec.Runner
	Inspector KindInspector
}

type StepError struct {
	Step    Step
	Subject string
	Err     error
}

type UnsafeTargetError struct {
	Path string
	Kind TargetKind
}
