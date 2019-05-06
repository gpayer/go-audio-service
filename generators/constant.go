package generators

import (
	"github.com/gpayer/go-audio-service/snd"
)

type Constant struct {
	samplerate uint32
	Value      float32
}

func NewConstant(samplerate uint32, v float32) *Constant {
	return &Constant{
		samplerate: samplerate,
		Value:      v,
	}
}

func (c *Constant) SetOutput(f snd.Writable) {
	f.SetReadable(c)
}

func (c *Constant) Read(samples *snd.Samples) {
	v := c.Value
	for i := 0; i < len(samples.Frames); i++ {
		samples.Frames[i].L = v
		samples.Frames[i].R = v
	}
}

func (c *Constant) ReadStateless(samples *snd.Samples, freq float32, _ *snd.NoteState) {
	c.Read(samples)
}
