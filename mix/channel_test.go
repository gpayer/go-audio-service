package mix

import (
	"github.com/gpayer/go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

type readTestData struct{}

func (r *readTestData) Read(samples *snd.Samples) int {
	samples.Frames[0] = snd.Sample{L: 1.0, R: 1.0}
	samples.Frames[1] = snd.Sample{L: -0.5, R: -0.5}
	return 2
}

func TestRead(t *testing.T) {
	assert := assert.New(t)
	samples := &snd.Samples{
		SampleRate: uint32(22000),
		Frames:     make([]snd.Sample, 2),
	}
	channel := NewChannel(22000)
	readable := &readableFunc{
		fn: func(samples *snd.Samples) {
			samples.Frames[0] = snd.Sample{L: 1.0, R: 1.0}
			samples.Frames[1] = snd.Sample{L: -0.5, R: -0.5}
		},
	}
	channel.SetReadable(readable)

	channel.Read(samples)

	assert.Equal(float32(1.0), samples.Frames[0].L)
	assert.Equal(float32(1.0), samples.Frames[0].R)

	channel.SetGain(0.5)
	channel.Read(samples)

	assert.Equal(float32(0.5), samples.Frames[0].L)
	assert.Equal(float32(0.5), samples.Frames[0].R)

	channel.SetGain(1.0)
	channel.SetPan(1.0)
	channel.Read(samples)

	assert.Equal(float32(1.0), samples.Frames[0].L)
	assert.Equal(float32(0.0), samples.Frames[0].R)
	assert.Equal(float32(-0.5), samples.Frames[1].L)
	assert.Equal(float32(0.0), samples.Frames[1].R)

	channel.SetPan(-1.0)
	channel.Read(samples)

	assert.Equal(float32(0.0), samples.Frames[0].L)
	assert.Equal(float32(1.0), samples.Frames[0].R)
	assert.Equal(float32(0.0), samples.Frames[1].L)
	assert.Equal(float32(-0.5), samples.Frames[1].R)
}
