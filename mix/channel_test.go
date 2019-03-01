package mix

import (
	"go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	assert := assert.New(t)
	samples := &snd.Samples{
		SampleRate: uint32(22000),
	}
	channel := NewChannel(22000)
	ch := make(chan *snd.Samples)
	channel.out = ch
	resetsamples := func() {
		samples.Frames = make([]snd.Sample, 0)
		samples.Add(snd.Sample{L: 1.0, R: 1.0}, snd.Sample{L: -0.5, R: -0.5})
	}
	resetsamples()

	writefunc := func() {
		assert.Nil(channel.Write(samples))
	}
	go writefunc()

	result := <-ch

	assert.Len(result.Frames, 2)
	assert.Equal(float32(1.0), result.Frames[0].L)
	assert.Equal(float32(1.0), result.Frames[0].R)

	resetsamples()
	channel.SetGain(0.5)
	go writefunc()
	result = <-ch

	assert.Equal(float32(0.5), result.Frames[0].L)
	assert.Equal(float32(0.5), result.Frames[0].R)

	resetsamples()
	channel.SetGain(1.0)
	channel.SetPan(1.0)
	go writefunc()
	result = <-ch

	assert.Equal(float32(1.0), result.Frames[0].L)
	assert.Equal(float32(0.0), result.Frames[0].R)
	assert.Equal(float32(-0.5), result.Frames[1].L)
	assert.Equal(float32(0.0), result.Frames[1].R)

	resetsamples()
	channel.SetPan(-1.0)
	go writefunc()
	result = <-ch

	assert.Equal(float32(0.0), result.Frames[0].L)
	assert.Equal(float32(1.0), result.Frames[0].R)
	assert.Equal(float32(0.0), result.Frames[1].L)
	assert.Equal(float32(-0.5), result.Frames[1].R)
}
