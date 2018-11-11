package logftext

import (
	"math/rand"
	"testing"

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
