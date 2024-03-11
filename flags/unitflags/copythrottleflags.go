package unitflags

import (
	"github.com/aggellos2001/go-robocopy/flags"
)

type UnitFlags flags.Flag

const (
	Kilobytes UnitFlags = 1 << iota
	Megabytes
	Gigabytes
)

func (c UnitFlags) String() string {
	if c == Kilobytes {
		return "k"
	} else if c == Megabytes {
		return "m"
	} else {
		return "g"
	}
}
