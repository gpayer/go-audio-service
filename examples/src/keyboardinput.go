package examples

import (
	"fmt"
	"go-audio-service/filters"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"image/color"

	"github.com/faiface/pixel/text"

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

type confSlider struct {
	w        float64
	txt      *text.Text
	slider   *Slider
	valueTxt *text.Text
	onchange func(v float32)
}

func newConfSlider(desc string, w, h float64, min, max, v float32, onchange func(v float32)) *confSlider {
	s := &confSlider{w: w, onchange: onchange}
	s.txt = text.New(pixel.ZV, FontService.Get("basic"))
	fmt.Fprint(s.txt, desc)
	s.slider = NewSlider(w, h, min, max, v)
	s.slider.OnChange(func(v float32) {
		s.onchange(v)
		s.valueTxt.Clear()
		fmt.Fprintf(s.valueTxt, "%.2f", s.slider.Value())
	})
	s.valueTxt = text.New(pixel.ZV, FontService.Get("basic"))
	s.valueTxt.Clear()
	fmt.Fprintf(s.valueTxt, "%.2f", s.slider.Value())
	return s
}

func (s *confSlider) Update(win *pixelgl.Window, dt float32, mat pixel.Matrix) {
	s.slider.Update(win, dt, mat.Moved(pixel.V(80, 0)))
	s.txt.Draw(win, mat.Moved(pixel.V(0, 10)))
	s.valueTxt.Draw(win, mat.Moved(pixel.V(120+s.w, 10)))
}

type keyboardExample struct {
	readable        snd.Readable
	instr           *DoubleOsci
	keys            map[pixelgl.Button]*keyDef
	keyIdx          []pixelgl.Button
	whiteKey        *imdraw.IMDraw
	blackKey        *imdraw.IMDraw
	sliderAttack    *confSlider
	sliderDecay     *confSlider
	sliderSustain   *confSlider
	sliderRelease   *confSlider
	sliderModFactor *confSlider
	sliderModGain   *confSlider
	whiteCanvas     *pixelgl.Canvas
	blackCanvas     *pixelgl.Canvas
}

func (k *keyboardExample) Init() {
	var attack, decay, sustain, release, modFactor, modGain float32
	attack = 0.05
	decay = 0.1
	sustain = 0.8
	release = 0.5
	modFactor = 2.3
	modGain = 0.1
	instr := NewDoubleOsci(attack, decay, sustain, release, modFactor, modGain)

	gain := filters.NewGain(0.3)
	gain.SetReadable(instr)
	compstate := filters.NewCompressorState()
	compstate.DefaultCompressor(44000)
	comp := filters.NewCompressor(44000, compstate)
	comp.SetReadable(gain)
	k.readable = comp
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

	k.sliderAttack = newConfSlider("Attack", 120, 30, 0.01, 1, attack, func(v float32) {
		k.instr.SetAttack(v)
	})
	k.sliderDecay = newConfSlider("Decay", 120, 30, 0, 3, decay, func(v float32) {
		k.instr.SetDecay(v)
	})
	k.sliderSustain = newConfSlider("Sustain", 120, 30, 0, 1, sustain, func(v float32) {
		k.instr.SetSustain(v)
	})
	k.sliderRelease = newConfSlider("Release", 120, 30, 0.01, 3, release, func(v float32) {
		k.instr.SetRelease(v)
	})
	k.sliderModFactor = newConfSlider("ModFactor", 120, 30, 0, 15, modFactor, func(v float32) {
		k.instr.SetModFactor(v)
	})
	k.sliderModGain = newConfSlider("ModGain", 120, 30, 0, 20, modGain, func(v float32) {
		k.instr.SetModGain(v)
	})
}

func (k *keyboardExample) Mounted() {
	GetOutput().SetReadable(k.readable)
	Start()
}

func (k *keyboardExample) Unmounted() {
	Stop()
}

func (k *keyboardExample) Update(win *pixelgl.Window, dt float32, mat pixel.Matrix) {
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

	if k.whiteCanvas == nil {
		k.whiteCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 50, 200))
		k.whiteKey.Draw(k.whiteCanvas)
		k.blackCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 30, 100))
		k.blackKey.Draw(k.blackCanvas)
	}

	bounds := win.Bounds()
	orig := mat.Moved(pixel.V(35, bounds.H()-70))

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
			k.whiteCanvas.DrawColorMask(win, orig.Moved(pixel.V(xWhite, 0)), maskcolor)
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
			k.blackCanvas.DrawColorMask(win, orig.Moved(pixel.V(xBlack, 25)), maskcolor)
		}
	}
	top := win.Bounds().H()
	k.sliderAttack.Update(win, dt, mat.Moved(pixel.V(20, top-280)))
	k.sliderDecay.Update(win, dt, mat.Moved(pixel.V(20, top-320)))
	k.sliderSustain.Update(win, dt, mat.Moved(pixel.V(20, top-360)))
	k.sliderRelease.Update(win, dt, mat.Moved(pixel.V(20, top-400)))
	k.sliderModFactor.Update(win, dt, mat.Moved(pixel.V(20, top-440)))
	k.sliderModGain.Update(win, dt, mat.Moved(pixel.V(20, top-480)))
}

func init() {
	AddExample("Keyboard", &keyboardExample{})
}
