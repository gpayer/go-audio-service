package main

import (
	"go-audio-service/generators"
	"go-audio-service/mix"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"time"
)

type noteShort struct {
	wait     time.Duration
	ch       int
	evtype   int
	notename string
	octave   int
	volume   float32
}

func createInstrument(a, d, s, r float32) *notes.NoteMultiplexer {
	rect1 := generators.NewRect(44000, 440.0)
	adsr1 := notes.NewAdsr(a, d, s, r)
	adsr1.SetReadable(rect1)
	multi1 := notes.NewNoteMultiplexer()
	multi1.SetReadable(adsr1)

	return multi1
}

func main() {
	output, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	piece := []noteShort{
		{0, 0, notes.Pressed, "G", 2, 1.0},
		{0, 1, notes.Pressed, "C", 3, 1.0},
		{100, 0, notes.Released, "C", 3, 0.0},
		{100, 0, notes.Released, "G", 2, 0.0},
		{100, 0, notes.Pressed, "E", 3, 1.0},
		{250, 0, notes.Released, "E", 3, 0.0},
		{0, 0, notes.Pressed, "G", 3, 1.0},
		{0, 1, notes.Pressed, "C", 2, 1.0},
		{500, 0, notes.Released, "G", 3, 0.0},
		{300, 1, notes.Released, "C", 2, 0.0},
	}

	var instr []*notes.NoteMultiplexer
	instr = append(instr, createInstrument(0.05, 0.3, 0.4, 0.8), createInstrument(0.1, 1.0, 0.0, 1.0))

	mixer := mix.NewMixer(44000)
	ch1 := mixer.GetChannel()
	ch1.SetReadable(instr[0])
	ch1.SetGain(0.3)

	ch2 := mixer.GetChannel()
	ch2.SetReadable(instr[1])
	ch2.SetGain(0.3)

	mixer.SetGain(0.6)
	output.SetReadable(mixer)

	err = output.Start()
	if err != nil {
		panic(err)
	}

	for _, n := range piece {
		if n.wait > 0 {
			time.Sleep(n.wait * time.Millisecond)
		}
		instr[n.ch].SendNoteEvent(notes.NewNoteEvent(n.evtype, notes.Note(n.notename, n.octave), n.volume))
	}
	time.Sleep(1000 * time.Millisecond)

	_ = output.Stop()
}
