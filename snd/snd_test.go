package snd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	assert := assert.New(t)
	samples := &Samples{SampleRate: 22000}
	s := Sample{L: 0.5, R: 0.2}
	s2 := Sample{L: 0.1, R: 0.1}
	samples.Add(s, s, s, s2)
	assert.Len(samples.Frames, 4)
	assert.Equal(float32(0.5), samples.Frames[1].L)
	assert.Equal(float32(0.1), samples.Frames[3].R)
}
