package filters

import (
	"fmt"
	"go-audio-service/snd"
)

type LowPassFilter struct {
	buffer []snd.Sample
	state  *BiquadState
	rate   uint32
	out    snd.Filter
}

func NewLowPass(rate uint32, cutoff, resonance float32) *LowPassFilter {
	state := &BiquadState{}
	state.LowPass(rate, cutoff, resonance)
	return &LowPassFilter{
		state: state,
		rate:  rate,
	}
}

func (f *LowPassFilter) Write(samples *snd.Samples) error {
	if samples.SampleRate != f.rate {
		return fmt.Errorf("incompatible sample rate: %d != %d", samples.SampleRate, f.rate)
	}
	if f.out != nil {
		f.state.Process(samples.Frames, samples.Frames)

		return f.out.Write(samples)
	}
	return nil
}

func (f *LowPassFilter) SetOutput(out snd.Filter) {
	f.out = out
}
