package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"time"
)

func runSin(output *snd.Output) error {
	sin := generators.NewSin(880)
	fminput, _ := sin.GetInput("fm")
	fmmod := generators.NewSin(880)
	fmgain := filters.NewGain(.2)
	fmgain.SetReadable(fmmod)
	fminput.SetReadable(fmgain)
	cont := notes.NewContinuousNote(notes.Note(notes.C, 4))
	cont.SetReadable(sin)
	output.SetReadable(cont)

	err := output.Start()
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	_ = output.Stop()
	return nil
}

func init() {
	AddExample("Sin", runSin)
}
