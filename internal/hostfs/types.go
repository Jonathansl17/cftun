package hostfs

type FS struct{}

type Kind int

const (
	KindMissing Kind = iota
	KindRegular
	KindDirectory
	KindSymlink
	KindOther
)

type Temp struct {
	Path string
}
