package examples

import (
	"flag"
	"fmt"
	"go-audio-service/filters"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"image/color"
	"pixelext/nodes"
	"pixelext/ui"

	"github.com/rakyll/portmidi"

	"github.com/faiface/pixel/text"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var midiDeviceID = flag.Int("devid", 0, "midi device ID")

type keyDef struct {
	note    notes.NoteValue
	pressed bool
	white   bool
}

type confSlider struct {
	nodes.BaseNode
	hbox     *ui.HBox
	w        float64
	txt      *nodes.Text
	slider   *ui.Slider
	valueTxt *nodes.Text
	onchange func(v float32)
}

func newConfSlider(desc string, w, h float64, min, max, v float32, onchange func(v float32)) *confSlider {
	s := &confSlider{
		BaseNode: *nodes.NewBaseNode(""),
		w:        w,
		onchange: onchange,
	}
	s.Self = s
	s.hbox = ui.NewHBox("hbox")
	s.AddChild(s.hbox)

	s.txt = nodes.NewText("txt", "basic")
	s.txt.Printf(desc)
	s.hbox.AddChild(s.txt)

	sliderval := nodes.NewBaseNode("")
	s.slider = ui.NewSlider("slider", min, max, v)
	s.slider.SetBounds(pixel.R(0, 0, w, h))
	s.slider.OnChange(func(v float32) {
		s.onchange(v)
		s.valueTxt.Clear()
		s.valueTxt.Printf("%.2f", s.slider.Value())
	})
	sliderval.AddChild(s.slider)

	s.valueTxt = nodes.NewText("valuetxt", "basic")
	s.valueTxt.SetZIndex(10)
	s.valueTxt.SetZeroAlignment(nodes.AlignmentCenter)
	s.valueTxt.SetPos(pixel.V(w/2, h/2))
	s.valueTxt.Printf("%.2f", v)
	sliderval.AddChild(s.valueTxt)
	s.hbox.AddChild(sliderval)

	return s
}

type keyboardExample struct {
	nodes.BaseNode
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
	txtOsci1        *text.Text
	toggleOsci1     *ValueToggle
	txtOsci2        *text.Text
	toggleOsci2     *ValueToggle
	portIn          *portmidi.Stream
}

func (k *keyboardExample) Init() {
	var attack, decay, sustain, release, modFactor, modGain float32
	attack = 0.05
	decay = 0.1
	sustain = 0.8
	release = 0.5
	modFactor = 2.3
	modGain = 0.1
	instr := NewDoubleOsci(attack, decay, sustain, release, modFactor, modGain, OsciSin, OsciSin)

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
	k.sliderAttack.SetPos(pixel.V(20, 320))
	k.AddChild(k.sliderAttack)

	k.sliderDecay = newConfSlider("Decay", 120, 30, 0, 3, decay, func(v float32) {
		k.instr.SetDecay(v)
	})
	k.sliderDecay.SetPos(pixel.V(20, 280))
	k.AddChild(k.sliderDecay)

	k.sliderSustain = newConfSlider("Sustain", 120, 30, 0, 1, sustain, func(v float32) {
		k.instr.SetSustain(v)
	})
	k.sliderSustain.SetPos(pixel.V(20, 240))
	k.AddChild(k.sliderSustain)

	k.sliderRelease = newConfSlider("Release", 120, 30, 0.01, 3, release, func(v float32) {
		k.instr.SetRelease(v)
	})
	k.sliderRelease.SetPos(pixel.V(20, 200))
	k.AddChild(k.sliderRelease)

	k.sliderModFactor = newConfSlider("ModFactor", 120, 30, 0, 15, modFactor, func(v float32) {
		k.instr.SetModFactor(v)
	})
	k.sliderModFactor.SetPos(pixel.V(20, 160))
	k.AddChild(k.sliderModFactor)

	k.sliderModGain = newConfSlider("ModGain", 120, 30, 0, 20, modGain, func(v float32) {
		k.instr.SetModGain(v)
	})
	k.sliderModGain.SetPos(pixel.V(20, 120))
	k.AddChild(k.sliderModGain)

	k.txtOsci1 = text.New(pixel.ZV, FontService.Get("basic"))
	fmt.Fprintf(k.txtOsci1, "OSCI1")
	k.toggleOsci1 = NewValueToggle(80, 15, func(val int) {
		k.instr.SetOsciType(1, OsciType(val))
	})
	k.toggleOsci1.AddValue("sin", int(OsciSin))
	k.toggleOsci1.AddValue("rect", int(OsciRect))
	k.txtOsci2 = text.New(pixel.ZV, FontService.Get("basic"))
	fmt.Fprintf(k.txtOsci2, "OSCI2")
	k.toggleOsci2 = NewValueToggle(80, 15, func(val int) {
		k.instr.SetOsciType(2, OsciType(val))
	})
	k.toggleOsci2.AddValue("sin", int(OsciSin))
	k.toggleOsci2.AddValue("rect", int(OsciRect))
}

func (k *keyboardExample) Mount() {
	if *midiDeviceID > 0 {
		portIn, err := portmidi.NewInputStream(portmidi.DeviceID(*midiDeviceID), 256)
		if err != nil {
			panic(err)
		}
		k.portIn = portIn
	}
	GetOutput().SetReadable(k.readable)
	Start()
}

func (k *keyboardExample) Unmount() {
	if *midiDeviceID > 0 && k.portIn != nil {
		k.portIn.Close()
		k.portIn = nil
	}
	Stop()
}

func (k *keyboardExample) Update(dt float64) {
	ev := nodes.Events()
	if ev.JustPressed(pixelgl.KeyQ) {
		SwitchScene("main")
	} else {
		for key, def := range k.keys {
			if ev.Pressed(key) && !k.keys[key].pressed {
				k.instr.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, def.note, 0.6))
				k.keys[key].pressed = true
			} else if !ev.Pressed(key) && k.keys[key].pressed {
				k.instr.SendNoteEvent(notes.NewNoteEvent(notes.Released, def.note, 0.0))
				k.keys[key].pressed = false
			}
		}
	}

	if *midiDeviceID > 0 && k.portIn != nil {
		midiEvents, err := k.portIn.Read(128)
		if err == nil {
			for _, ev := range midiEvents {
				if ev.Status == 0x90 {
					k.instr.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, notes.MidiToNote(ev.Data1), float32(ev.Data2)/127.0))
				} else if ev.Status == 0x80 {
					k.instr.SendNoteEvent(notes.NewNoteEvent(notes.Released, notes.MidiToNote(ev.Data1), 0.0))
				}
			}
		}
	}
}

func (k *keyboardExample) Draw(win *pixelgl.Window, mat pixel.Matrix) {
	if k.whiteCanvas == nil {
		k.whiteCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 50, 200))
		k.whiteKey.Draw(k.whiteCanvas)
		k.blackCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 30, 100))
		k.blackKey.Draw(k.blackCanvas)
	}

	//bounds := win.Bounds()
	orig := mat.Moved(pixel.V(35, 530))

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
	/*k.sliderAttack.Update(win, dt, mat.Moved(pixel.V(20, top-280)))
	k.sliderDecay.Update(win, dt, mat.Moved(pixel.V(20, top-320)))
	k.sliderSustain.Update(win, dt, mat.Moved(pixel.V(20, top-360)))
	k.sliderRelease.Update(win, dt, mat.Moved(pixel.V(20, top-400)))
	k.sliderModFactor.Update(win, dt, mat.Moved(pixel.V(20, top-440)))
	k.sliderModGain.Update(win, dt, mat.Moved(pixel.V(20, top-480)))

	k.txtOsci1.Draw(win, mat.Moved(pixel.V(400, top-260)))
	k.toggleOsci1.Update(win, dt, mat.Moved(pixel.V(450, top-260)))
	k.txtOsci2.Draw(win, mat.Moved(pixel.V(400, top-280)))
	k.toggleOsci2.Update(win, dt, mat.Moved(pixel.V(450, top-280)))*/
}

func init() {
	k := &keyboardExample{
		BaseNode: *nodes.NewBaseNode("keyboard"),
	}
	k.Self = k
	AddExample("Keyboard", k)
}
