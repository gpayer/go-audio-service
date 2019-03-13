package main

import (
	"go-audio-service/filters"
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

	sin := generators.NewSin(880)
	fminput, _ := sin.GetInput("fm")
	fmmod := generators.NewSin(880)
	fmgain := filters.NewGain(.2)
	fmgain.SetReadable(fmmod)
	fminput.SetReadable(fmgain)
	cont := notes.NewContinuousNote(notes.Note(notes.C, 4))
	cont.SetReadable(sin)
	output.SetReadable(cont)

	err = output.Start()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	_ = output.Stop()
}
