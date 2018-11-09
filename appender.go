package logftext

import (
	"io"
	"os"

	"github.com/ssgreg/logf"
)

// NewAppender returns a new logf.WriteAppender with the given Writer and
// EncoderConfig.
//
// NewAppender is safe to use for colored logs.
func NewAppender(w io.Writer, cfg EncoderConfig) logf.Appender {
	if f, ok := w.(*os.File); ok {
		err := enableSeqTTY(f.Fd(), true)

		if cfg.NoColor == nil {
			noColor := false
			if err != nil {
				noColor = true
			}
			if checkNoColor() {
				noColor = true
			}
			cfg.NoColor = &noColor
		}
	}

	return logf.NewWriteAppender(w, NewEncoder(cfg))
}

// checkNoColor checks for NO_COLORS environment variable to disable color
// output.
//
// All command-line software which outputs text with ANSI color added should
// check for the presence of a NO_COLOR environment variable that, when
// present (regardless of its value), prevents the addition of ANSI color.
func checkNoColor() bool {
	_, ok := os.LookupEnv("NO_COLOR")

	return ok
}
