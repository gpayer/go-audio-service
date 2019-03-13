package generators

import (
	"go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSin(t *testing.T) {
	assert := assert.New(t)

	notestate := &snd.NoteState{
		Timecode: 0,
		On:       true,
	}

	sin := NewSin(10)
	samples := snd.NewSamples(1000, 100)
	sin.ReadStateless(samples, 0, notestate)

	assert.Equal(float32(0.0), samples.Frames[0].L)
}
