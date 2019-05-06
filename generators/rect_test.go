package generators

import (
	"github.com/gpayer/go-audio-service/filters"
	"github.com/gpayer/go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRectStateless(t *testing.T) {
	assert := assert.New(t)
	r := NewRect(0, 1)
	samples := snd.NewSamples(1000, 100)

	var high float32 = 0.5
	var low float32 = -0.5
	notestate := &snd.NoteState{Timecode: 0, On: true}

	r.ReadStateless(samples, 20.0, notestate)

	assert.Equal(high, samples.Frames[0].L)
	assert.Equal(high, samples.Frames[51].L)
	assert.Equal(low, samples.Frames[26].L)
	assert.Equal(low, samples.Frames[76].L)

	notestate.Timecode = 26
	r.ReadStateless(samples, 20.0, notestate)
	assert.Equal(low, samples.Frames[0].L)
	assert.Equal(low, samples.Frames[51].L)
	assert.Equal(high, samples.Frames[26].L)
	assert.Equal(high, samples.Frames[76].L)
}

func TestRectAM(t *testing.T) {
	assert := assert.New(t)
	r := NewRect(0, 20)
	samples := snd.NewSamples(1000, 100)

	ammod := NewConstant(0, 0.1)
	am, ok := r.GetInput("am")
	assert.True(ok)
	am.SetReadable(ammod)

	r.ReadStateless(samples, 0.0, snd.EmptyNoteState)

	var high float32 = 0.55
	var low float32 = -0.55
	assert.Equal(high, samples.Frames[0].L)
	assert.Equal(high, samples.Frames[51].L)
	assert.Equal(low, samples.Frames[26].L)
	assert.Equal(low, samples.Frames[76].L)

	ammodsin := NewSin(10)
	gain := filters.NewGain(2)
	gain.SetReadable(ammodsin)
	am.SetReadable(gain)

	r.ReadStateless(samples, 0.0, snd.EmptyNoteState)

	assert.Equal(float32(0.5), samples.Frames[0].L)
	assert.True(float32(0.5)-samples.Frames[50].L < 0.001)
	assert.Equal(float32(-1), samples.Frames[25].L)
}
