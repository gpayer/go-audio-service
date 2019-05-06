package main

import (
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/notes"
	"github.com/gpayer/go-audio-service/snd"
	"time"
)

func main() {
	output, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	rect := generators.NewRect(44000, 440)
	cont := notes.NewContinuousNote(notes.Note(notes.C, 3))
	cont.SetReadable(rect)
	output.SetReadable(cont)

	err = output.Start()
	if err != nil {
		panic(err)
	}

	var notelist []notes.NoteValue = []notes.NoteValue{
		notes.Note(notes.D, 3),
		notes.Note(notes.E, 3),
		notes.Note(notes.F, 3),
		notes.Note(notes.G, 3),
		notes.Note(notes.A, 3),
		notes.Note(notes.H, 3),
		notes.Note(notes.C, 4),
	}
	for _, nv := range notelist {
		time.Sleep(250 * time.Millisecond)
		cont.SetNote(nv)
	}
	time.Sleep(250 * time.Millisecond)

	_ = output.Stop()
}
