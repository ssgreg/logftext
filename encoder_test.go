package logftext

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/ssgreg/logf"
	"github.com/stretchr/testify/require"
)

type encoderTestCase struct {
	Name    string
	Entry   []logf.Entry
	Golden  string
	NoColor bool
	Config  EncoderConfig
}

type MyInt int
type MyUint uint
type MyBool bool
type MyFloat float64
type MyString string

type user struct {
	Name string `json:"name"`
}
type users []*user

func (u *user) EncodeLogfObject(enc logf.FieldEncoder) error {
	enc.EncodeFieldString("name", u.Name)

	return nil
}

func (u users) EncodeLogfArray(enc logf.TypeEncoder) error {
	for i := range u {
		enc.EncodeTypeObject(u[i])
	}

	return nil
}

func TestEncoder(t *testing.T) {
	testCases := []encoderTestCase{
		{
			"OnlyMessage",
			[]logf.Entry{
				{
					LoggerID: int32(rand.Int()),
					Level:    logf.LevelInfo,
					Text:     "message",
				},
			},
			`Jan  1 00:00:00.000 |INFO| message` + "\n",
			true,
			EncoderConfig{},
		},
		{
			"CustomTimeEncoder",
			[]logf.Entry{
				{
					LoggerID: int32(rand.Int()),
					Level:    logf.LevelWarn,
					Text:     "another message",
				},
			},
			`0001-01-01T00:00:000 |WARN| another message` + "\n",
			true,
			EncoderConfig{
				EncodeTime: logf.RFC3339TimeEncoder,
			},
		},
		{
			"WithArrayAndObjectFields",
			[]logf.Entry{
				{
					LoggerID: int32(rand.Int()),
					Level:    logf.LevelWarn,
					Text:     "message",
					Fields: []logf.Field{
						logf.Object("o", &user{"n"}),
						logf.Array("a", users{{"n1"}, {"n2"}}),
					},
				},
			},
			`Jan  1 00:00:00.000 |WARN| message o={"name":"n"} a=[{"name":"n1"},{"name":"n2"}]` + "\n",
			true,
			EncoderConfig{},
		},
		{
			"WithFieldsAndDerivedFields",
			[]logf.Entry{
				{
					LoggerID: int32(rand.Int()),
					Level:    logf.LevelDebug,
					Text:     "message",
					Fields: []logf.Field{
						logf.ConstInts("ints", []int{0, 1}),
					},
					DerivedFields: []logf.Field{
						logf.ConstBytes("bytes", []byte("!")),
					},
				},
			},
			`Jan  1 00:00:00.000 |DEBU| message bytes="IQ==" ints=[0,1]` + "\n",
			true,
			EncoderConfig{},
		},
		{
			"WithFieldsAndDerivedFieldsAndCaller",
			[]logf.Entry{
				{
					LoggerID: int32(rand.Int()),
					Level:    logf.LevelWarn,
					Text:     "message",
					Fields: []logf.Field{
						logf.ConstFloats32("fts", []float32{0.1, 9}),
					},
					DerivedFields: []logf.Field{
						logf.ConstDurations("drs", []time.Duration{time.Second}),
					},
					Caller: logf.EntryCaller{
						PC:        0,
						File:      "/a/b/c/f.go",
						Line:      6,
						Specified: true,
					},
				},
			},
			`Jan  1 00:00:00.000 |WARN| message drs=["1s"] fts=[0.1,9] @"c/f.go:6"` + "\n",
			true,
			EncoderConfig{},
		},
		{
			"WithCompleteListOfAnyFields",
			[]logf.Entry{
				{
					LoggerID: int32(rand.Int()),
					Level:    logf.LevelWarn,
					Text:     "message",
					Fields: []logf.Field{
						logf.Any("s", "3"),
						logf.Any("b", true),
						logf.Any("i", int(1)),
						logf.Any("i8", int8(1)),
						logf.Any("i16", int16(1)),
						logf.Any("i32", int32(1)),
						logf.Any("i64", int64(1)),
						logf.Any("u", uint(1)),
						logf.Any("u8", uint8(1)),
						logf.Any("u16", uint16(1)),
						logf.Any("u32", uint32(1)),
						logf.Any("u64", uint64(1)),
						logf.Any("f32", float32(1)),
						logf.Any("f64", float64(1)),
						logf.Any("e", errors.New("1")),
						logf.Any("ri", MyInt(1)),
						logf.Any("ru", MyUint(1)),
						logf.Any("rb", MyBool(true)),
						logf.Any("rs", MyString("2")),
						logf.Any("rf", MyFloat(1)),
					},
				},
			},
			`Jan  1 00:00:00.000 |WARN| message s="3" b=true i=1 i8=1 i16=1 i32=1 i64=1 u=1 u8=1 u16=1 u32=1 u64=1 f32=1 f64=1 e="1" ri=1 ru=1 rb=true rs="2" rf=1` + "\n",
			true,
			EncoderConfig{},
		},
		{
			"WithCompleteListOfFields",
			[]logf.Entry{
				{
					LoggerID: int32(rand.Int()),
					Level:    logf.LevelError,
					Text:     "message",
					Fields: []logf.Field{
						logf.String("s", "1"),
						logf.Bool("b", true),
						logf.Int("i", 1),
						logf.Int8("i8", 1),
						logf.Int16("i16", 1),
						logf.Int32("i32", 1),
						logf.Int64("i64", 1),
						logf.Uint("u", 1),
						logf.Uint8("u8", 1),
						logf.Uint16("u16", 1),
						logf.Uint32("u32", 1),
						logf.Uint64("u64", 1),
						logf.Float32("f32", 1),
						logf.Float64("f64", 1),
						logf.Duration("d", 1),
						logf.ConstBools("bs", []bool{true, false}),
						logf.ConstInts("is", []int{0, 1}),
						logf.ConstInts8("is8", []int8{0, 1}),
						logf.ConstInts16("is16", []int16{0, 1}),
						logf.ConstInts32("is32", []int32{0, 1}),
						logf.ConstInts64("is64", []int64{0, 1}),
						logf.ConstUints("us", []uint{0, 1}),
						logf.ConstUints8("us8", []uint8{0, 1}),
						logf.ConstUints16("us16", []uint16{0, 1}),
						logf.ConstUints32("us32", []uint32{0, 1}),
						logf.ConstUints64("us64", []uint64{0, 1}),
						logf.ConstFloats32("fs32", []float32{0, 1}),
						logf.ConstFloats64("fs64", []float64{0, 1}),
					},
				},
			},
			`Jan  1 00:00:00.000 |ERRO| message s="1" b=true i=1 i8=1 i16=1 i32=1 i64=1 u=1 u8=1 u16=1 u32=1 u64=1 f32=1 f64=1 d="1ns" bs=[true,false] is=[0,1] is8=[0,1] is16=[0,1] is32=[0,1] is64=[0,1] us=[0,1] us8=[0,1] us16=[0,1] us32=[0,1] us64=[0,1] fs32=[0,1] fs64=[0,1]` + "\n",
			true,
			EncoderConfig{},
		},
		{
			"WithFieldsAndCallerAndName",
			[]logf.Entry{
				{
					LoggerID:   int32(rand.Int()),
					Level:      logf.LevelWarn,
					Text:       "message",
					LoggerName: "name",
					Fields: []logf.Field{
						logf.String("test", "f"),
					},
					Caller: logf.EntryCaller{
						PC:        0,
						File:      "/a/b/c/f.go",
						Line:      6,
						Specified: true,
					},
				},
			},
			`Jan  1 00:00:00.000 |WARN| name: message test="f" @"c/f.go:6"` + "\n",
			true,
			EncoderConfig{},
		},
		{
			"WithFieldsAndCallerAndNameAndColored",
			[]logf.Entry{
				{
					LoggerID:   int32(rand.Int()),
					Level:      logf.LevelWarn,
					Text:       "message",
					LoggerName: "name",
					Fields: []logf.Field{
						logf.String("test", "f"),
					},
					Caller: logf.EntryCaller{
						PC:        0,
						File:      "/a/b/c/f.go",
						Line:      6,
						Specified: true,
					},
				},
			},
			"\x1b[90mJan  1 00:00:00.000\x1b[0m |\x1b[93;7mWARN\x1b[0m| \x1b[90mname:\x1b[0m \x1b[97mmessage\x1b[0m \x1b[32mtest\x1b[0m\x1b[90m=\x1b[0m\"f\"\x1b[90m @\"c/f.go:6\"\x1b[0m" + "\n",
			false,
			EncoderConfig{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			b := logf.NewBuffer()

			for _, e := range tc.Entry {
				cfg := tc.Config
				cfg.NoColor = &tc.NoColor

				enc := NewEncoder(cfg)
				enc.Encode(b, e)
			}

			require.EqualValues(t, tc.Golden, b.String())
		})
	}
}
