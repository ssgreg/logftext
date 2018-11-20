package logftext

import (
	"os"
)

// EnableSeqTTY enables possibility to use escape sequences in TTY if possible.
func EnableSeqTTY(f *os.File, flag bool) bool {
	return enableSeqTTY(f.Fd(), flag) == nil
}

// CheckNoColor checks for NO_COLORS environment variable to disable color
// output.
//
// All command-line software which outputs text with ANSI color added should
// check for the presence of a NO_COLOR environment variable that, when
// present (regardless of its value), prevents the addition of ANSI color.
func CheckNoColor() bool {
	_, ok := os.LookupEnv("NO_COLOR")

	return ok
}
