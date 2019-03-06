package mix

import (
	"go-audio-service/snd"
	"math"
)

// Channel is a special volume and pan filter, which can be connected to a mixer
type Channel struct {
	gain       float32
	pan        float32
	samplerate uint32
	readable   snd.Readable
}

// NewChannel creates a new Channel instance
func NewChannel(samplerate uint32) *Channel {
	return &Channel{
		samplerate: samplerate,
		gain:       1.0,
		pan:        0.0,
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

func (ch *Channel) Read(samples *snd.Samples) int {
	length := ch.readable.Read(samples)
	scale := float32(1.0 - math.Abs(float64(ch.pan))*.5)
	lgain := (ch.pan + 1.0) * scale
	rgain := (1.0 - ch.pan) * scale
	for idx, sample := range samples.Frames {
		sample.L *= lgain * ch.gain
		sample.R *= rgain * ch.gain
		samples.Frames[idx] = sample
	}
	return length
}
