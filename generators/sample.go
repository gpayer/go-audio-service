package generators

import (
	"go-audio-service/snd"
)

type Sample struct {
	samples *snd.Samples
}

func NewSample(samples *snd.Samples) *Sample {
	return &Sample{
		samples: samples,
	}
}

func (s *Sample) Read(samples *snd.Samples) {
	s.ReadStateless(samples, 0, snd.EmptyNoteState)
}

func (s *Sample) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	if state.On {
		pos := state.Timecode
		srclen := uint32(len(s.samples.Frames))
		for i := 0; i < len(samples.Frames); i++ {
			if pos < srclen {
				samples.Frames[i] = s.samples.Frames[pos]
			} else {
				samples.Frames[i] = snd.Sample{L: 0, R: 0}
			}
			pos++
		}
	}
}
