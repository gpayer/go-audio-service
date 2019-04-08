package generators

import (
	"go-audio-service/snd"
)

type Sample struct {
	samples   *snd.Samples
	playFull  bool
	noteEnded bool
}

func NewSample(samples *snd.Samples) *Sample {
	return &Sample{
		samples:   samples,
		playFull:  true,
		noteEnded: false,
	}
}

func (s *Sample) SetPlayFull(playFull bool) {
	s.playFull = playFull
}

func (s *Sample) Read(samples *snd.Samples) {
	s.ReadStateless(samples, 0, snd.EmptyNoteState)
}

func (s *Sample) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	if !s.playFull && !state.On {
		return
	}
	pos := state.Timecode
	srclen := uint32(len(s.samples.Frames))
	for i := 0; i < len(samples.Frames); i++ {
		if pos < srclen {
			samples.Frames[i].L = s.samples.Frames[pos].L * state.Volume
			samples.Frames[i].R = s.samples.Frames[pos].R * state.Volume
			s.noteEnded = false
		} else {
			samples.Frames[i] = snd.Sample{L: 0, R: 0}
			s.noteEnded = true
		}
		pos++
	}
}

func (s *Sample) NoteEnded() bool {
	return s.noteEnded
}
