package dcopyflags

import (
	"strings"

	"github.com/aggellos2001/go-robocopy/flags"
)

type DCopyFlags flags.Flag

const (
	D       DCopyFlags = 1 << iota // Data
	A                              // Attributes
	T                              // Time stamps
	E                              // Extended attribute
	X                              // Skip alt data streams
	Default = D | A                // Data, Attributes
)

func (flag DCopyFlags) String() string {
	var r strings.Builder
	if flags.Has(flag, D) {
		r.WriteString("D")
	}
	if flags.Has(flag, A) {
		r.WriteString("A")
	}
	if flags.Has(flag, T) {
		r.WriteString("T")
	}
	if flags.Has(flag, E) {
		r.WriteString("E")
	}
	if flags.Has(flag, X) {
		r.WriteString("X")
	}
	return r.String()
}
