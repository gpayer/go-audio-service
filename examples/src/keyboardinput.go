package examples

import (
	"go-audio-service/filters"
	"go-audio-service/notes"
	"go-audio-service/snd"

	"github.com/faiface/pixel/pixelgl"
)

type keyboardExample struct {
	readable snd.Readable
	instr    *notes.NoteMultiplexer
	keys     map[pixelgl.Button]notes.NoteValue
}

func (k *keyboardExample) Init() {
	instr := createInstrument(1, 0.01, 0.1, 0.8, 0.5)

	gain := filters.NewGain(0.3)
	gain.SetReadable(instr)
	k.readable = gain
	k.instr = instr

	k.keys = make(map[pixelgl.Button]notes.NoteValue, 16)
	k.keys[pixelgl.KeyA] = notes.Note(notes.C, 4)
	k.keys[pixelgl.KeyS] = notes.Note(notes.D, 4)
	k.keys[pixelgl.KeyD] = notes.Note(notes.E, 4)
	k.keys[pixelgl.KeyF] = notes.Note(notes.F, 4)
	k.keys[pixelgl.KeyG] = notes.Note(notes.G, 4)
	k.keys[pixelgl.KeyH] = notes.Note(notes.A, 4)
	k.keys[pixelgl.KeyJ] = notes.Note(notes.H, 4)
	k.keys[pixelgl.KeyK] = notes.Note(notes.C, 5)
}

func (k *keyboardExample) Mounted() {
	GetOutput().SetReadable(k.readable)
	Start()
}

func (k *keyboardExample) Unmounted() {
	Stop()
}

func (k *keyboardExample) Update(win *pixelgl.Window, dt float32) {
	if win.JustPressed(pixelgl.KeyQ) {
		SwitchScene("main")
	} else {
		for key, note := range k.keys {
			if win.JustPressed(key) {
				k.instr.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, note, 0.6))
			} else if win.JustReleased(key) {
				k.instr.SendNoteEvent(notes.NewNoteEvent(notes.Released, note, 0.0))
			}
		}
	}
}

func init() {
	AddExample("Keyboard", &keyboardExample{})
}
