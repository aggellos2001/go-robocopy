package aflags

import (
	"strings"

	"github.com/aggellos2001/go-robocopy/flags"
)

type AFlags flags.Flag

const (
	R AFlags = 1 << iota // Read only
	A                    // Archive
	S                    // System
	H                    // Hidden
	C                    // Compressed
	N                    // Not content indexed
	E                    // Encrypted
	T                    // Temporary
	O                    // Offline - only used in the [/a-]: flag
)

func (flag AFlags) String() string {
	var r strings.Builder
	if flags.Has(flag, R) {
		r.WriteString("R")
	}
	if flags.Has(flag, A) {
		r.WriteString("A")
	}
	if flags.Has(flag, S) {
		r.WriteString("S")
	}
	if flags.Has(flag, H) {
		r.WriteString("H")
	}
	if flags.Has(flag, C) {
		r.WriteString("C")
	}
	if flags.Has(flag, N) {
		r.WriteString("N")
	}
	if flags.Has(flag, E) {
		r.WriteString("E")
	}
	if flags.Has(flag, T) {
		r.WriteString("T")
	}
	if flags.Has(flag, O) {
		r.WriteString("O")
	}
	return r.String()
}
