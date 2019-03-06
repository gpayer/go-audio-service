package mix

import (
	"go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

type readableFunc struct {
	fn func(samples *snd.Samples)
}

func (r *readableFunc) Read(samples *snd.Samples) int {
	r.fn(samples)
	return len(samples.Frames)
}

func TestMixer(t *testing.T) {
	assert := assert.New(t)
	samples := &snd.Samples{
		SampleRate: 22000,
		Frames:     make([]snd.Sample, 5),
	}
	readable1 := &readableFunc{
		fn: func(samples *snd.Samples) {
			for i := 0; i < 5; i++ {
				samples.Frames[i].L = 0.3
				samples.Frames[i].R = 0.3
			}
		},
	}
	readable2 := &readableFunc{
		fn: func(samples *snd.Samples) {
			for i := 0; i < 5; i++ {
				samples.Frames[i].L = 0.2
				samples.Frames[i].R = 0.2
			}
		},
	}

	m := NewMixer(22000)
	ch1 := m.GetChannel()
	ch1.SetReadable(readable1)
	ch2 := m.GetChannel()
	ch2.SetReadable(readable2)

	length := m.Read(samples)

	assert.True(m.running)
	m.Stop()
	assert.False(m.running)
	assert.Equal(length, 5)
	assert.Equal(float32(0.5), samples.Frames[0].L)
	assert.Equal(float32(0.5), samples.Frames[0].R)
}
