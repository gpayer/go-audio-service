package generators

import (
	"go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRectStateless(t *testing.T) {
	assert := assert.New(t)
	r := NewRect(0, 1)
	samples := snd.NewSamples(1000, 100)

	var high float32 = 0.5
	var low float32 = -0.5

	r.ReadStateless(samples, 20.0, 0)

	assert.Equal(high, samples.Frames[0].L)
	assert.Equal(high, samples.Frames[51].L)
	assert.Equal(low, samples.Frames[26].L)
	assert.Equal(low, samples.Frames[76].L)

	r.ReadStateless(samples, 20.0, 26)
	assert.Equal(low, samples.Frames[0].L)
	assert.Equal(low, samples.Frames[51].L)
	assert.Equal(high, samples.Frames[26].L)
	assert.Equal(high, samples.Frames[76].L)
}
