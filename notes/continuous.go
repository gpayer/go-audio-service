package notes

import (
	"go-audio-service/snd"
)

type ContinuousNote struct {
	r     snd.Readable
	freq  float32
	state *snd.NoteState
}

func NewContinuousNote(note NoteValue) *ContinuousNote {
	return &ContinuousNote{
		freq:  float32(note),
		state: &snd.NoteState{Volume: 1.0, Timecode: 0, On: true},
	}
}

func (c *ContinuousNote) SetReadable(r snd.Readable) {
	c.r = r
}

func (c *ContinuousNote) Read(samples *snd.Samples) {
	length := len(samples.Frames)
	c.r.ReadStateless(samples, c.freq, c.state)
	c.state.Timecode += uint32(length)
	//c.state.Timecode %= samples.SampleRate
}

func (c *ContinuousNote) ReadStateless(samples *snd.Samples, freq float32, _ *snd.NoteState) {
	c.Read(samples)
}

func (c *ContinuousNote) SetNote(note NoteValue) {
	c.freq = float32(note)
	c.state.Timecode = 0
}
