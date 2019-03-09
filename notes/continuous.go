package notes

import (
	"go-audio-service/snd"
)

type ContinuousNote struct {
	r        snd.Readable
	freq     float32
	timecode uint32
}

func NewContinuousNote(note NoteValue) *ContinuousNote {
	return &ContinuousNote{
		freq:     float32(note),
		timecode: 0,
	}
}

func (c *ContinuousNote) SetReadable(r snd.Readable) {
	c.r = r
}

func (c *ContinuousNote) Read(samples *snd.Samples) {
	length := len(samples.Frames)
	c.r.ReadStateless(samples, c.freq, c.timecode, true)
	c.timecode += uint32(length)
	c.timecode %= samples.SampleRate
}

func (c *ContinuousNote) ReadStateless(samples *snd.Samples, freq float32, timecode uint32, _ bool) {
	c.Read(samples)
}

func (c *ContinuousNote) SetNote(note NoteValue) {
	c.freq = float32(note)
}
