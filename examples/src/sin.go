package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func runSin(output snd.IOutput, _ *pixelgl.Window) error {
	sin := generators.NewSin(880)

	fminput, _ := sin.GetInput("fm")
	fmmod := generators.NewSin(880)
	fmgain := filters.NewGain(.9)
	fmgain.SetReadable(fmmod)
	fminput.SetReadable(fmgain)

	aminput, _ := sin.GetInput("am")
	ammod := generators.NewSin(3)
	amgain := filters.NewGain(.2)
	amgain.SetReadable(ammod)
	aminput.SetReadable(amgain)

	cont := notes.NewContinuousNote(notes.Note(notes.C, 3))
	cont.SetReadable(sin)

	gain := filters.NewGain(.3)
	gain.SetReadable(cont)
	output.SetReadable(gain)

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
