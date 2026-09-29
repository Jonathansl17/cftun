package cli

type FixedLocator struct {
	Path string
}

func (l FixedLocator) ConfigPath() string {
	return l.Path
}
