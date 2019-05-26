package mix

import (
	"math"

	"github.com/gpayer/go-audio-service/snd"
)

// Channel is a special volume and pan filter, which can be connected to a mixer
type Channel struct {
	gain       float32
	pan        float32
	samplerate uint32
	enabled    bool
	readable   snd.Readable
}

// NewChannel creates a new Channel instance
func NewChannel(samplerate uint32) *Channel {
	return &Channel{
		samplerate: samplerate,
		gain:       1.0,
		pan:        0.0,
		enabled:    true,
	}
}

func (ch *Channel) SetReadable(r snd.Readable) {
	ch.readable = r
}

func clamp(v float32, min float32, max float32) float32 {
	if v > max {
		return max
	} else if v < min {
		return min
	}
	return v
}

func (ch *Channel) SetEnabled(enabled bool) {
	ch.enabled = enabled
}

func (ch *Channel) Enabled() bool {
	return ch.enabled
}

// SetGain sets gain value
func (ch *Channel) SetGain(gain float32) {
	ch.gain = clamp(gain, 0.0, 1.0)
}

// Gain returns gain value
func (ch *Channel) Gain() float32 {
	return ch.gain
}

// SetPan sets pan value
func (ch *Channel) SetPan(pan float32) {
	ch.pan = clamp(pan, -1.0, 1.0)
}

// Pan returns pan value
func (ch *Channel) Pan() float32 {
	return ch.pan
}

func (ch *Channel) Read(samples *snd.Samples) {
	ch.ReadStateless(samples, 0, snd.EmptyNoteState)
}

func (ch *Channel) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	if !ch.enabled {
		for i := 0; i < len(samples.Frames); i++ {
			samples.Frames[i] = snd.Sample{
				L: 0,
				R: 0,
			}
		}
	}
	ch.readable.ReadStateless(samples, freq, state)
	scale := float32(1.0 - math.Abs(float64(ch.pan))*.5)
	lgain := (ch.pan + 1.0) * scale
	rgain := (1.0 - ch.pan) * scale
	for idx, sample := range samples.Frames {
		sample.L *= lgain * ch.gain
		sample.R *= rgain * ch.gain
		samples.Frames[idx] = sample
	}
}
