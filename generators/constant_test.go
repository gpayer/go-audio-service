package generators

import (
	"go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstant(t *testing.T) {
	assert := assert.New(t)
	c := NewConstant(44000, 30.0)
	samples := &snd.Samples{
		SampleRate: 44000,
		Frames:     make([]snd.Sample, 128),
	}

	c.Read(samples)

	for _, fr := range samples.Frames {
		assert.Equal(float32(30.0), fr.L)
		assert.Equal(float32(30.0), fr.R)
	}
}

func TestConstantStateless(t *testing.T) {
	assert := assert.New(t)
	c := NewConstant(44000, 30.0)
	samples := &snd.Samples{
		SampleRate: 44000,
		Frames:     make([]snd.Sample, 128),
	}

	c.ReadStateless(samples, 220.0, 0, true)

	for _, fr := range samples.Frames {
		assert.Equal(float32(30.0), fr.L)
		assert.Equal(float32(30.0), fr.R)
	}
}
