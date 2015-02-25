package terraform

//go:generate stringer -type=DiffSource diff_source.go

// DiffSource is a bitmask type that can say where a value is set.
type DiffSource int

const (
	DiffSourceInvalid DiffSource = 0
	DiffSourceConfig  DiffSource = 1 << iota
	DiffSourceState
)
