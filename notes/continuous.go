package notes

import (
	"go-audio-service/generators"
	"go-audio-service/snd"
)

type ContinuousNote struct {
	g        generators.Generator
	freq     float32
	timecode uint32
}

func NewContinuousNote(note NoteValue) *ContinuousNote {
	return &ContinuousNote{
		freq:     float32(note),
		timecode: 0,
	}
}

func (c *ContinuousNote) SetGenerator(g generators.Generator) {
	c.g = g
}

func (c *ContinuousNote) Read(samples *snd.Samples) int {
	length := len(samples.Frames)
	c.g.ReadStateless(samples, c.freq, c.timecode)
	c.timecode += uint32(length)
	c.timecode %= samples.SampleRate
	return length
}
