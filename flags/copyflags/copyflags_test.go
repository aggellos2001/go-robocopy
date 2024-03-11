package copyflags

import "testing"

func TestString(t *testing.T) {
	var have []CopyFlags = []CopyFlags{
		D, A, T, X, S, O, U, Default,
		D | A, D | T, D | X, D | S, D | O, D | U,
	}
	var want []string = []string{
		"D", "A", "T", "X", "S", "O", "U", "DAT",
		"DA", "DT", "DX", "DS", "DO", "DU",
	}
	for i, flag := range have {
		if flag.String() != want[i] {
			t.Errorf("want: %s, have: %s", want[i], flag.String())
		}
	}
}
