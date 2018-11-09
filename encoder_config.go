package logftext

import (
	"time"

	"github.com/ssgreg/logf"
)

// EncoderConfig allows to configure text Encoder.
type EncoderConfig struct {
	DisableFieldName   bool
	DisableFieldCaller bool

	EncodeTime     logf.TimeEncoder
	EncodeDuration logf.DurationEncoder
	EncodeError    logf.ErrorEncoder
	EncodeCaller   logf.CallerEncoder
}

// WithDefaults returns the new config in which all uninitialized fields are
// filled with their default values.
func (c EncoderConfig) WithDefaults() EncoderConfig {
	// Handle defaults for type encoder.
	if c.EncodeDuration == nil {
		c.EncodeDuration = logf.StringDurationEncoder
	}
	if c.EncodeTime == nil {
		c.EncodeTime = logf.LayoutTimeEncoder(time.StampMilli)
	}
	if c.EncodeError == nil {
		c.EncodeError = logf.DefaultErrorEncoder
	}
	if c.EncodeCaller == nil {
		c.EncodeCaller = logf.ShortCallerEncoder
	}

	return c
}
