package dcopyflags

import "testing"

func TestString(t *testing.T) {
	var have []DCopyFlags = []DCopyFlags{
		D, A, T, E, X, Default,
		D | A, D | T, D | E, D | X,
	}
	var want []string = []string{
		"D", "A", "T", "E", "X", "DA",
		"DA", "DT", "DE", "DX",
	}
	for i, flag := range have {
		if flag.String() != want[i] {
			t.Errorf("want: %s, have: %s", want[i], flag.String())
		}
	}
}
