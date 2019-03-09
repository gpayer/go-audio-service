package main

import (
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"time"
)

func main() {
	output, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	rect := generators.NewRect(44000, 440.0)
	multi := notes.NewNoteMultiplexer()
	multi.SetReadable(rect)
	output.SetReadable(multi)

	err = output.Start()
	if err != nil {
		panic(err)
	}

	multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, notes.Note(notes.C, 3), 0.1))
	time.Sleep(500 * time.Millisecond)
	multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, notes.Note(notes.E, 3), 0.1))
	time.Sleep(500 * time.Millisecond)
	multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, notes.Note(notes.G, 3), 0.1))
	time.Sleep(750 * time.Millisecond)
	multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, notes.Note(notes.G, 2), 0.1))
	time.Sleep(1000 * time.Millisecond)
	multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, notes.Note(notes.C, 3), 0.0))
	multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, notes.Note(notes.E, 3), 0.0))
	multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, notes.Note(notes.G, 3), 0.0))
	time.Sleep(500 * time.Millisecond)
	multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, notes.Note(notes.G, 2), 0.0))

	_ = output.Stop()
}
