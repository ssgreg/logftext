package logftext

import (
	"math/rand"
	"testing"

	"github.com/ssgreg/logf"
	"github.com/stretchr/testify/require"
)

type encoderTestCase struct {
	Name   string
	Entry  []logf.Entry
	Golden string
}

func TestEncoder(t *testing.T) {
	testCases := []encoderTestCase{
		{
			"Simple",
			[]logf.Entry{
				{
					LoggerID: int32(rand.Int()),
					Level:    logf.LevelInfo,
					Text:     "message",
				},
			},
			`Jan  1 00:00:00.000 |INFO| message` + "\n",
		},
	}

	noColor := true
	enc := NewEncoder(EncoderConfig{
		NoColor: &noColor,
	})

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			b := logf.NewBuffer()
			for _, e := range tc.Entry {
				enc.Encode(b, e)
			}

			require.EqualValues(t, tc.Golden, b.String())
		})
	}

}
