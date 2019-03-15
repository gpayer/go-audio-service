package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"time"
)

func run(output snd.IOutput) error {
	rect := generators.NewRect(44000, 440)

	fm, _ := rect.GetInput("fm")
	fmmod := generators.NewSin(200)
	fmgain := filters.NewGain(10)
	fmgain.SetReadable(fmmod)
	fm.SetReadable(fmgain)

	am, _ := rect.GetInput("am")
	ammod := generators.NewSin(3)
	amgain := filters.NewGain(.3)
	amgain.SetReadable(ammod)
	am.SetReadable(amgain)

	cont := notes.NewContinuousNote(notes.Note(notes.C, 3))
	cont.SetReadable(rect)

	gain := filters.NewGain(0.3)
	gain.SetReadable(cont)
	output.SetReadable(gain)

	Start()
	time.Sleep(time.Second)
	Stop()

	return nil
}

func init() {
	AddExample("Rect", run)
}
