package size

import "github.com/dustin/go-humanize"

// IEC Sizes.
// kibis of bits
const (
	B = 1 << (iota * 10)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB
)

// Format returns a string representation of a size in a uint64 number of bytes.
// It uses the IEC standard (kiB, MiB, GiB, TiB)
func Format(s uint64) string {
	return humanize.IBytes(s)
}

// ParseBytes parses a string representation of bytes into the number
// of bytes it represents.
func ParseBytes(s string) (uint64, error) {
	return humanize.ParseBytes(s)
}
