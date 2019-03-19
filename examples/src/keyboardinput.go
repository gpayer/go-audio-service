package examples

import (
	"go-audio-service/filters"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type keyDef struct {
	note    notes.NoteValue
	pressed bool
	white   bool
}

type keyboardExample struct {
	readable snd.Readable
	instr    *notes.NoteMultiplexer
	keys     map[pixelgl.Button]*keyDef
	keyIdx   []pixelgl.Button
	whiteKey *imdraw.IMDraw
	blackKey *imdraw.IMDraw
}

func (k *keyboardExample) Init() {
	instr := createInstrument(1, 0.01, 0.1, 0.8, 0.5)

	gain := filters.NewGain(0.3)
	gain.SetReadable(instr)
	k.readable = gain
	k.instr = instr

	k.keyIdx = []pixelgl.Button{
		pixelgl.KeyA, pixelgl.KeyW,
		pixelgl.KeyS, pixelgl.KeyE,
		pixelgl.KeyD, pixelgl.KeyF,
		pixelgl.KeyT, pixelgl.KeyG,
		pixelgl.KeyY, pixelgl.KeyH,
		pixelgl.KeyU, pixelgl.KeyJ,
		pixelgl.KeyK,
	}
	k.keys = make(map[pixelgl.Button]*keyDef, 16)
	k.keys[pixelgl.KeyA] = &keyDef{notes.Note(notes.C, 4), false, true}
	k.keys[pixelgl.KeyW] = &keyDef{notes.Note(notes.Csharp, 4), false, false}
	k.keys[pixelgl.KeyS] = &keyDef{notes.Note(notes.D, 4), false, true}
	k.keys[pixelgl.KeyE] = &keyDef{notes.Note(notes.Dsharp, 4), false, false}
	k.keys[pixelgl.KeyD] = &keyDef{notes.Note(notes.E, 4), false, true}
	k.keys[pixelgl.KeyF] = &keyDef{notes.Note(notes.F, 4), false, true}
	k.keys[pixelgl.KeyT] = &keyDef{notes.Note(notes.Fsharp, 4), false, false}
	k.keys[pixelgl.KeyG] = &keyDef{notes.Note(notes.G, 4), false, true}
	k.keys[pixelgl.KeyY] = &keyDef{notes.Note(notes.Gsharp, 4), false, false}
	k.keys[pixelgl.KeyH] = &keyDef{notes.Note(notes.A, 4), false, true}
	k.keys[pixelgl.KeyU] = &keyDef{notes.Note(notes.Hb, 4), false, false}
	k.keys[pixelgl.KeyJ] = &keyDef{notes.Note(notes.H, 4), false, true}
	k.keys[pixelgl.KeyK] = &keyDef{notes.Note(notes.C, 5), false, true}

	k.whiteKey = imdraw.New(nil)
	k.whiteKey.Color = colornames.White
	k.whiteKey.Push(pixel.V(0, 0), pixel.V(0, 200), pixel.V(50, 200), pixel.V(50, 0))
	k.whiteKey.Polygon(0)

	k.blackKey = imdraw.New(nil)
	k.blackKey.Color = colornames.White
	k.blackKey.Push(pixel.V(0, 0), pixel.V(0, 100), pixel.V(30, 100), pixel.V(30, 0))
	k.blackKey.Polygon(0)
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
		for key, def := range k.keys {
			if win.Pressed(key) && !k.keys[key].pressed {
				k.instr.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, def.note, 0.6))
				k.keys[key].pressed = true
			} else if !win.Pressed(key) && k.keys[key].pressed {
				k.instr.SendNoteEvent(notes.NewNoteEvent(notes.Released, def.note, 0.0))
				k.keys[key].pressed = false
			}
		}
	}

	whiteKey := pixelgl.NewCanvas(pixel.R(0, 0, 50, 200))
	k.whiteKey.Draw(whiteKey)
	blackKey := pixelgl.NewCanvas(pixel.R(0, 0, 30, 100))
	k.blackKey.Draw(blackKey)

	bounds := win.Bounds()
	orig := pixel.IM.Moved(pixel.V(35, bounds.H()-60))

	xWhite := 0.0
	xBlack := 27.0 - 54.0
	for _, v := range k.keyIdx {
		def := k.keys[v]
		if def.white {
			var maskcolor color.RGBA
			if def.pressed {
				maskcolor = colornames.Yellow
			} else {
				maskcolor = colornames.White
			}
			whiteKey.DrawColorMask(win, orig.Moved(pixel.V(xWhite, 0)), maskcolor)
			xWhite += 54
		}
	}
	for _, v := range k.keyIdx {
		def := k.keys[v]
		if def.white {
			xBlack += 54
		} else {
			var maskcolor color.RGBA
			if def.pressed {
				maskcolor = colornames.Yellow
			} else {
				maskcolor = colornames.Black
			}
			blackKey.DrawColorMask(win, orig.Moved(pixel.V(xBlack, 25)), maskcolor)
		}
	}
}

func init() {
	AddExample("Keyboard", &keyboardExample{})
}
