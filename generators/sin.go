package generators

import (
	"go-audio-service/snd"
	"math"
)

type Sin struct {
	snd.BasicWritableProvider
	freq          float32
	samplerate    uint32
	dphi          float32
	fm            *snd.BasicConnector
	am            *snd.BasicConnector
	FreqModFactor float32
}

func NewSin(freq float32) *Sin {
	s := &Sin{freq: freq}
	s.InitBasicWritableProvider()
	s.fm = s.AddInput("fm", 0.0)
	s.am = s.AddInput("am", .5)
	s.FreqModFactor = 0
	return s
}

func (s *Sin) Read(samples *snd.Samples) {
	if s.samplerate != samples.SampleRate {
		s.samplerate = samples.SampleRate
		s.dphi = float32(2.0*math.Pi) / (float32(samples.SampleRate) / s.freq)
	}
	s.ReadStateless(samples, s.freq, snd.EmptyNoteState)
}

func (s *Sin) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	if freq != s.freq && freq > 0 || s.samplerate != samples.SampleRate {
		s.samplerate = samples.SampleRate
		if freq > 0 {
			s.freq = freq
		}
		s.dphi = float32(2.0*math.Pi) / (float32(samples.SampleRate) / s.freq)
	}
	phi := float32(state.Timecode) * s.dphi

	fm := s.fm.ReadBuffered(samples.SampleRate, len(samples.Frames), freq*s.FreqModFactor, state)
	am := s.am.ReadBuffered(samples.SampleRate, len(samples.Frames), 0, state)

	for i := 0; i < len(samples.Frames); i++ {
		v := float32(math.Sin(float64(phi+fm.Frames[i].L))) * (1.0 + am.Frames[i].L)
		samples.Frames[i].L = v
		samples.Frames[i].R = v
		phi += s.dphi
	}
}
