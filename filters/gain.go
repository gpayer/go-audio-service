package filters

import (
	"github.com/gpayer/go-audio-service/snd"
)

type Gain struct {
	gain     float32
	readable snd.Readable
}

func NewGain(gain float32) *Gain {
	return &Gain{gain: gain}
}

func (g *Gain) SetGain(gain float32) {
	g.gain = gain
}

func (g *Gain) SetReadable(r snd.Readable) {
	g.readable = r
}

func (g *Gain) process(samples *snd.Samples) {
	for i := 0; i < len(samples.Frames); i++ {
		samples.Frames[i].L *= g.gain
		samples.Frames[i].R *= g.gain
	}
}

func (g *Gain) Read(samples *snd.Samples) {
	g.readable.Read(samples)
	g.process(samples)
}

func (g *Gain) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	g.readable.ReadStateless(samples, freq, state)
	g.process(samples)
}
