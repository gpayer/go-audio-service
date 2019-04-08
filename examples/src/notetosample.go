package examples

import (
	"fmt"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"pixelext/services"
)

type NoteToSample struct {
	samples   map[float32]*generators.Sample
	noteEnded bool
}

func NewNoteToSample() *NoteToSample {
	return &NoteToSample{
		samples:   make(map[float32]*generators.Sample, 0),
		noteEnded: false,
	}
}

func (n *NoteToSample) AddSample(note notes.NoteValue, path string) {
	s, err := services.ResourceManager().LoadSample(path)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		n.samples[float32(note)] = generators.NewSample(s)
	}
}

func (n *NoteToSample) Read(samples *snd.Samples) {
}

func (n *NoteToSample) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	s, ok := n.samples[freq]
	if ok {
		s.ReadStateless(samples, freq, state)
		n.noteEnded = s.NoteEnded()
	}
}

func (n *NoteToSample) NoteEnded() bool {
	return n.noteEnded
}
