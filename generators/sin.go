package generators

import (
	"go-audio-service/snd"
	"math"
)

type Sin struct {
	freq float32
	dphi float32
}

func NewSin(freq float32) *Sin {
	return &Sin{freq: freq}
}

func (s *Sin) Read(samples *snd.Samples) {
	s.ReadStateless(samples, s.freq, snd.EmptyNoteState)
}

func (s *Sin) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	if freq != s.freq {
		s.freq = freq
		s.dphi = float32(2.0*math.Pi) / (float32(samples.SampleRate) / freq)
	}
	phi := float32(state.Timecode) * s.dphi

	for i := 0; i < len(samples.Frames); i++ {
		v := float32(math.Cos(float64(phi)))
		samples.Frames[i].L = v
		samples.Frames[i].R = v
		phi += s.dphi
	}
}
