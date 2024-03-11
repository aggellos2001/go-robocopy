package copyflags

import (
	"strings"

	"github.com/aggellos2001/go-robocopy/flags"
)

type CopyFlags flags.Flag

const (
	D       CopyFlags   = 1 << iota // Data
	A                               // Attributes
	T                               // Time stamps
	X                               // Skip alt data streams
	S                               // NTFS access control list (ACL)
	O                               // Owner information
	U                               // Auditing information
	Default = D | A | T             // Data, Attributes, Time stamps
)

func (flag CopyFlags) String() string {
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
	if flags.Has(flag, X) {
		r.WriteString("X")
	}
	if flags.Has(flag, S) {
		r.WriteString("S")
	}
	if flags.Has(flag, O) {
		r.WriteString("O")
	}
	if flags.Has(flag, U) {
		r.WriteString("U")
	}
	return r.String()
}
