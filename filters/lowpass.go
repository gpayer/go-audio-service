package filters

import (
	"go-audio-service/snd"
)

type LowPassFilter struct {
	snd.BasicWritableProvider
	state        *BiquadState
	rate         uint32
	cutoff       float32
	resonance    float32
	cutoffInput  *snd.BasicConnector
	readable     snd.Readable
	cutoffValues *snd.Samples
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
	lowpass.InitBasicWritableProvider()
	lowpass.cutoffInput = lowpass.AddInput("cutoff")
	lowpass.cutoffValues = snd.NewSamples(rate, 512)
	return lowpass
}

func (f *LowPassFilter) Read(samples *snd.Samples) {
	f.ReadStateless(samples, 0, 0, true)
}

func (f *LowPassFilter) ReadStateless(samples *snd.Samples, freq float32, timecode uint32, on bool) {
	if len(samples.Frames) != len(f.cutoffValues.Frames) {
		f.cutoffValues = snd.NewSamples(samples.SampleRate, len(samples.Frames))
	}
	f.cutoffInput.ReadStateless(f.cutoffValues, freq, timecode, on)
	if f.cutoffValues.Valid {
		newCutoff := f.cutoffValues.Frames[0].L
		if newCutoff != f.cutoff {
			f.state.LowPass(f.rate, newCutoff, f.resonance)
			f.cutoff = newCutoff
		}
	}
	f.readable.ReadStateless(samples, freq, timecode, on)
	f.state.Process(samples.Frames, samples.Frames)
}

func (f *LowPassFilter) SetReadable(r snd.Readable) {
	f.readable = r
}

func (f *LowPassFilter) SetOutput(out snd.Writable) {
	out.SetReadable(f)
}
