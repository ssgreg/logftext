package logftext

import (
	"time"

	"github.com/ssgreg/logf"
)

// NewEncoder creates the new instance of the text Encoder with the given
// EncoderConfig. It delegates field encoding to the internal json Encoder.
//
// It's a caller responsibility to handle colored output in a prover way.
// The best choice here is to use NewAppender function instead of of a
// creation of Encoder by Yourselves.
var NewEncoder = encoderGetter(
	func(cfg EncoderConfig) logf.Encoder {
		cfg = cfg.WithDefaults()

		return &encoder{
			cfg,
			logf.NewJSONTypeEncoderFactory(logf.JSONEncoderConfig{
				EncodeTime:     cfg.EncodeTime,
				EncodeDuration: cfg.EncodeDuration,
				EncodeError:    cfg.EncodeError,
			}),
			nil,
			logf.NewCache(100),
			0,
			EscapeSequence{*cfg.NoColor},
		}
	},
)

type encoderGetter func(cfg EncoderConfig) logf.Encoder

func (c encoderGetter) Default() logf.Encoder {
	return c(EncoderConfig{})
}

type encoder struct {
	EncoderConfig
	mf logf.TypeEncoderFactory

	buf         *logf.Buffer
	cache       *logf.Cache
	startBufLen int

	eseq EscapeSequence
}

func (f *encoder) Encode(buf *logf.Buffer, e logf.Entry) error {
	// TODO: move to clone
	f.buf = buf
	f.startBufLen = f.buf.Len()

	// Time.
	f.eseq.At(f.buf, EscBrightBlack, func() {
		appendTime(e.Time, f.buf, f.EncodeTime, f.mf.TypeEncoder(buf))
	})

	// Level.
	f.appendSeparator()
	appendLevel(buf, f.eseq, e.Level)

	// Logger name.
	if !f.DisableFieldName && e.LoggerName != "" {
		f.appendSeparator()
		f.eseq.At(f.buf, EscBrightBlack, func() {
			f.buf.AppendString(e.LoggerName)
			f.buf.AppendByte(':')
		})
	}

	// Message.
	f.appendSeparator()
	f.eseq.At(f.buf, EscBrightWhite, func() {
		f.buf.AppendString(e.Text)
	})

	// Logger's fields.
	if bytes, ok := f.cache.Get(e.LoggerID); ok {
		buf.AppendBytes(bytes)
	} else {
		le := buf.Len()
		for _, field := range e.DerivedFields {
			field.Accept(f)
		}

		bf := make([]byte, buf.Len()-le)
		copy(bf, buf.Data[le:])
		f.cache.Set(e.LoggerID, bf)
	}

	// Entry's fields.
	for _, field := range e.Fields {
		field.Accept(f)
	}

	// Caller.
	if !f.DisableFieldCaller && e.Caller.Specified {
		f.eseq.At(f.buf, EscBrightBlack, func() {
			f.appendSeparator()
			f.buf.AppendByte('@')
			f.EncodeCaller(e.Caller, f.mf.TypeEncoder(f.buf))
		})
	}

	buf.AppendByte('\n')

	return nil
}

func (f *encoder) EncodeFieldAny(k string, v interface{}) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeAny(v)
}

func (f *encoder) EncodeFieldBool(k string, v bool) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeBool(v)
}

func (f *encoder) EncodeFieldInt64(k string, v int64) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeInt64(v)
}

func (f *encoder) EncodeFieldInt32(k string, v int32) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeInt32(v)
}

func (f *encoder) EncodeFieldInt16(k string, v int16) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeInt16(v)
}

func (f *encoder) EncodeFieldInt8(k string, v int8) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeInt8(v)
}

func (f *encoder) EncodeFieldUint64(k string, v uint64) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeUint64(v)
}

func (f *encoder) EncodeFieldUint32(k string, v uint32) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeUint32(v)
}

func (f *encoder) EncodeFieldUint16(k string, v uint16) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeUint16(v)
}

func (f *encoder) EncodeFieldUint8(k string, v uint8) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeUint8(v)
}

func (f *encoder) EncodeFieldFloat64(k string, v float64) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeFloat64(v)
}

func (f *encoder) EncodeFieldFloat32(k string, v float32) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeFloat32(v)
}

func (f *encoder) EncodeFieldString(k string, v string) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeString(v)
}

func (f *encoder) EncodeFieldDuration(k string, v time.Duration) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeDuration(v)
}

func (f *encoder) EncodeFieldError(k string, v error) {
	f.EncodeError(k, v, f)
}

func (f *encoder) EncodeFieldTime(k string, v time.Time) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeTime(v)
}

func (f *encoder) EncodeFieldArray(k string, v logf.ArrayEncoder) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeArray(v)
}

func (f *encoder) EncodeFieldObject(k string, v logf.ObjectEncoder) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeObject(v)
}

func (f *encoder) EncodeFieldBytes(k string, v []byte) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeBytes(v)
}

func (f *encoder) EncodeFieldBools(k string, v []bool) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeBools(v)
}

func (f *encoder) EncodeFieldInts64(k string, v []int64) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeInts64(v)
}

func (f *encoder) EncodeFieldInts32(k string, v []int32) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeInts32(v)
}

func (f *encoder) EncodeFieldInts16(k string, v []int16) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeInts16(v)
}

func (f *encoder) EncodeFieldInts8(k string, v []int8) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeInts8(v)
}

func (f *encoder) EncodeFieldUints64(k string, v []uint64) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeUints64(v)
}

func (f *encoder) EncodeFieldUints32(k string, v []uint32) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeUints32(v)
}

func (f *encoder) EncodeFieldUints16(k string, v []uint16) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeUints16(v)
}

func (f *encoder) EncodeFieldUints8(k string, v []uint8) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeUints8(v)
}

func (f *encoder) EncodeFieldFloats64(k string, v []float64) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeFloats64(v)
}

func (f *encoder) EncodeFieldFloats32(k string, v []float32) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeFloats32(v)
}

func (f *encoder) EncodeFieldDurations(k string, v []time.Duration) {
	f.addKey(k)
	f.mf.TypeEncoder(f.buf).EncodeTypeDurations(v)
}

func (f *encoder) appendSeparator() {
	if f.empty() {
		return
	}

	switch f.buf.Back() {
	case '=':
		return
	}
	f.buf.AppendByte(' ')
}

func (f *encoder) empty() bool {
	return f.buf.Len() == f.startBufLen
}

func (f *encoder) addKey(k string) {
	f.appendSeparator()
	f.eseq.At(f.buf, EscGreen, func() {
		f.buf.AppendString(k)
	})

	f.eseq.At(f.buf, EscBrightBlack, func() {
		f.buf.AppendByte('=')
	})
}

func appendLevel(buf *logf.Buffer, eseq EscapeSequence, lvl logf.Level) {
	buf.AppendByte('|')

	switch lvl {
	case logf.LevelDebug:
		eseq.At(buf, EscMagenta, func() {
			buf.AppendString("DEBU")
		})
	case logf.LevelInfo:
		eseq.At(buf, EscCyan, func() {
			buf.AppendString("INFO")
		})
	case logf.LevelWarn:
		eseq.At2(buf, EscBrightYellow, EscReverse, func() {
			buf.AppendString("WARN")
		})
	case logf.LevelError:
		eseq.At2(buf, EscBrightRed, EscReverse, func() {
			buf.AppendString("ERRO")
		})
	default:
		eseq.At(buf, EscBrightRed, func() {
			buf.AppendString("UNKN")
		})
	}

	buf.AppendByte('|')
}

func appendTime(t time.Time, buf *logf.Buffer, enc logf.TimeEncoder, encType logf.TypeEncoder) {
	start := buf.Len()
	enc(t, encType)
	end := buf.Len()

	// Get rid of possible quotes.
	if end != start {
		if buf.Data[start] == '"' && buf.Back() == '"' {
			copy(buf.Data[start:end], buf.Data[start+1:end-2])
			buf.Data = buf.Data[0 : end-2]
		}
	}
}
