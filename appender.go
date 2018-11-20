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
		ok = EnableSeqTTY(f, true)

		if cfg.NoColor == nil {
			// Set NoColor mode if Writer is not TTY or NO_COLOR environment
			// variable is set.
			noColor := !ok || CheckNoColor()
			cfg.NoColor = &noColor
		}
	}

	return logf.NewWriteAppender(w, NewEncoder(cfg))
}
