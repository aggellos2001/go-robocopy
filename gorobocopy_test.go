package gorobocopy

import (
	"slices"
	"testing"

	"github.com/aggellos2001/go-robocopy/flags/copyflags"
)

func TestCommandOne(t *testing.T) {

	cmd := NewRobocopy(
		"C:\\source",
		"D:\\destination",
		"*.*",
	)
	cmd.SetCopyOptions(&CopyOptions{
		E:    true,
		Mt:   4,
		Copy: copyflags.D | copyflags.A | copyflags.T,
	})

	cmd.SetFileSelectionOptions(&FileSelectionOptions{
		Xf: []string{"*.tmp", "*.bak"},
	})

	cmd.SetLoggingOptions(&LoggingOptions{
		V: true,
	})

	cmd.SetJobOptions(&JobOptions{
		Save: "MyJob",
		Quit: true,
	})

	want := []string{
		"C:\\source",
		"D:\\destination",
		"*.*",
		"/e",
		"/copy:DAT",
		"/mt:4",
		"/xf",
		"*.tmp",
		"*.bak",
		"/v",
		"/quit",
		"/save:MyJob",
	}
	have := cmd.GetCommandArgs()

	if !slices.Equal(want, have) {
		t.Errorf("have: %v\n,want: %v\n", have, want)
	}

}
