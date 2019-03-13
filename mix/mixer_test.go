package mix

import (
	"fmt"
	"go-audio-service/generators"
	"go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

type readableFunc struct {
	fn func(samples *snd.Samples)
}

func (r *readableFunc) Read(samples *snd.Samples) {
	r.fn(samples)
}

func (r *readableFunc) ReadStateless(samples *snd.Samples, freq float32, _ *snd.NoteState) {
	r.fn(samples)
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

	m.Read(samples)

	assert.True(m.running)
	m.Stop()
	assert.False(m.running)
	assert.Equal(float32(0.5), samples.Frames[0].L)
	assert.Equal(float32(0.5), samples.Frames[0].R)
}

func runMixerBenchmark(n int, b *testing.B) {
	m := NewMixer(44000)
	var rects []*generators.Rect
	for i := 0; i < n; i++ {
		r := generators.NewRect(44000, float32(800+n*10))
		rects = append(rects, r)
		ch := m.GetChannel()
		ch.SetReadable(r)
		ch.SetGain(0.5 / float32(n))
	}
	samples := snd.NewSamples(44000, 223)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Read(samples)
	}
}

func BenchmarkMixer(b *testing.B) {
	n := 1
	for n < 50 {
		b.Run(fmt.Sprintf("BenchmarkMixer%d", n), func(b *testing.B) {
			runMixerBenchmark(1, b)
		})
		n += 5
	}
}
