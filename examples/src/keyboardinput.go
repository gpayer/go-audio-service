package examples

import (
	"go-audio-service/filters"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func runKeyboardInput(output snd.IOutput, win *pixelgl.Window) error {
	instr := createInstrument(2, 0.5, 0.1, 0.8, 0.8)

	gain := filters.NewGain(0.3)
	gain.SetReadable(instr)
	output.SetReadable(gain)

	_ = output.Start()

	for {
		if win.JustPressed(pixelgl.KeyQ) {
			_ = output.Stop()
			return nil
		} else if win.JustPressed(pixelgl.KeyA) {
			instr.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, notes.Note(notes.C, 3), 0.2))
		} else if win.JustReleased(pixelgl.KeyA) {
			instr.SendNoteEvent(notes.NewNoteEvent(notes.Released, notes.Note(notes.C, 3), 0.0))
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func init() {
	AddExample("Keyboard", runKeyboardInput)
}
