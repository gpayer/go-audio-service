package filters

import (
	"github.com/gpayer/go-audio-service/snd"
)

type LowPassFilter struct {
	snd.BasicWritableProvider
	state       *BiquadState
	rate        uint32
	cutoff      float32
	resonance   float32
	cutoffInput *snd.BasicConnector
	readable    snd.Readable
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
	lowpass.cutoffInput = lowpass.AddInput("cutoff", cutoff)
	return lowpass
}

func (f *LowPassFilter) Read(samples *snd.Samples) {
	f.ReadStateless(samples, 0, snd.EmptyNoteState)
}

func (f *LowPassFilter) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	cutoffValues := f.cutoffInput.ReadBuffered(samples.SampleRate, len(samples.Frames), freq, state)
	newCutoff := cutoffValues.Frames[0].L
	if newCutoff != f.cutoff {
		f.state.LowPass(f.rate, newCutoff, f.resonance)
		f.cutoff = newCutoff
	}
	f.readable.ReadStateless(samples, freq, state)
	f.state.Process(samples.Frames, samples.Frames)
}

func (f *LowPassFilter) SetReadable(r snd.Readable) {
	f.readable = r
}

func (f *LowPassFilter) SetOutput(out snd.Writable) {
	out.SetReadable(f)
}
