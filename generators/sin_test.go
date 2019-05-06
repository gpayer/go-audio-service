package generators

import (
	"github.com/gpayer/go-audio-service/snd"
	"math"
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
	assert.True(samples.Frames[25].L > 0.49)
	assert.True(math.Abs(float64(samples.Frames[50].L)) < 0.001)
	assert.True(samples.Frames[75].L < -0.49)
}
