package filters

import (
	"fmt"
	"go-audio-service/snd"
)

type LowPassFilter struct {
	snd.BufferedInputProvider
	state       *BiquadState
	rate        uint32
	cutoff      float32
	resonance   float32
	cutoffInput *snd.InputBuffer
	out         snd.Input
}

func NewLowPass(rate uint32, cutoff, resonance float32) *LowPassFilter {
	state := &BiquadState{}
	state.LowPass(rate, cutoff, resonance)
	lowpass := &LowPassFilter{
		state:     state,
		rate:      rate,
		cutoff:    cutoff,
		resonance: resonance,
	}
	state.Reset()
	lowpass.Init()
	lowpass.cutoffInput = lowpass.AddInput("cutoff", rate)
	return lowpass
}

func (f *LowPassFilter) Write(samples *snd.Samples) error {
	if samples.SampleRate != f.rate {
		return fmt.Errorf("incompatible sample rate: %d != %d", samples.SampleRate, f.rate)
	}
	if f.out != nil {
		cutoffValues := f.cutoffInput.Read(len(samples.Frames))
		if len(cutoffValues.Frames) > 0 {
			newCutoff := cutoffValues.Frames[0].L
			if newCutoff != f.cutoff {
				f.state.LowPass(f.rate, newCutoff, f.resonance)
				f.cutoff = newCutoff
			}
		}
		f.state.Process(samples.Frames, samples.Frames)

		return f.out.Write(samples)
	}
	return nil
}

func (f *LowPassFilter) SetOutput(out snd.Input) {
	f.out = out
}
