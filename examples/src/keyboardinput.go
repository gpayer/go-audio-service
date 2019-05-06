package examples

import (
	"flag"
	"fmt"
	"image/color"
	"strconv"

	"github.com/gpayer/go-audio-service/filters"
	"github.com/gpayer/go-audio-service/notes"
	"github.com/gpayer/go-audio-service/snd"
	"github.com/gpayer/pixelext/nodes"
	"github.com/gpayer/pixelext/ui"

	"github.com/rakyll/portmidi"

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

func newConfSlider(grid *ui.Grid, desc string, w, h float64, min, max, v float32, onchange func(v float32)) {
	txt := ui.NewText("txt", "basic")
	txt.Printf(desc)
	grid.AddChild(txt)

	valueTxt := ui.NewText("valuetxt", "basic")

	sliderval := ui.NewUIBase("")
	sliderval.SetSize(pixel.V(w, h))
	slider := ui.NewSlider("slider", min, max, v)
	slider.SetSize(pixel.V(w, h))
	slider.OnChange(func(v float32) {
		onchange(v)
		valueTxt.Clear()
		valueTxt.Printf("%.2f", slider.Value())
	})
	sliderval.AddChild(slider)

	valueTxt.SetZIndex(10)
	valueTxt.SetAlignment(nodes.AlignmentCenter)
	valueTxt.SetPos(pixel.ZV)
	valueTxt.Printf("%.2f", v)
	sliderval.AddChild(valueTxt)
	grid.AddChild(sliderval)
}

type keyboardExample struct {
	nodes.BaseNode
	readable    snd.Readable
	instr       *DoubleOsci
	keys        map[pixelgl.Button]*keyDef
	keyIdx      []pixelgl.Button
	whiteKey    *imdraw.IMDraw
	blackKey    *imdraw.IMDraw
	whiteCanvas *pixelgl.Canvas
	blackCanvas *pixelgl.Canvas
	waveOsci1   *ui.ButtonGroup
	waveOsci2   *ui.ButtonGroup
	portIn      *portmidi.Stream
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

	grid := ui.NewGrid("grid", 2)
	grid.SetAlignment(nodes.AlignmentTopLeft)
	grid.SetPos(pixel.V(20, 320))
	k.AddChild(grid)

	newConfSlider(grid, "Attack", 120, 30, 0.01, 1, attack, func(v float32) {
		k.instr.SetAttack(v)
	})

	newConfSlider(grid, "Decay", 120, 30, 0, 3, decay, func(v float32) {
		k.instr.SetDecay(v)
	})

	newConfSlider(grid, "Sustain", 120, 30, 0, 1, sustain, func(v float32) {
		k.instr.SetSustain(v)
	})

	newConfSlider(grid, "Release", 120, 30, 0.01, 3, release, func(v float32) {
		k.instr.SetRelease(v)
	})

	newConfSlider(grid, "ModFactor", 120, 30, 0, 15, modFactor, func(v float32) {
		k.instr.SetModFactor(v)
	})

	newConfSlider(grid, "ModGain", 120, 30, 0, 20, modGain, func(v float32) {
		k.instr.SetModGain(v)
	})

	txt := ui.NewText("txtosci1", "basic")
	txt.Printf("OSCI1")
	txt.SetPos(pixel.V(400, 340))
	txt.SetAlignment(nodes.AlignmentCenterLeft)
	k.AddChild(txt)
	k.waveOsci1 = ui.NewButtonGroup("waveosci1", 0)
	k.waveOsci1.AddButton("sin", "sin", 0)
	k.waveOsci1.AddButton("rect", "rect", 0)
	k.waveOsci1.OnChange(func(v string) {
		var oscitype OsciType
		switch v {
		case "rect":
			oscitype = OsciRect
		default:
			oscitype = OsciSin
		}
		k.instr.SetOsciType(1, oscitype)
	})
	k.waveOsci1.SetPos(pixel.V(450, 340))
	k.waveOsci1.SetAlignment(nodes.AlignmentCenterLeft)
	k.AddChild(k.waveOsci1)

	txt = ui.NewText("txtosci2", "basic")
	txt.Printf("OSCI2")
	txt.SetPos(pixel.V(400, 280))
	txt.SetAlignment(nodes.AlignmentCenterLeft)
	k.AddChild(txt)
	k.waveOsci2 = ui.NewButtonGroup("waveosci2", 0)
	k.waveOsci2.AddButton("sin", "sin", 0)
	k.waveOsci2.AddButton("rect", "rect", 0)
	k.waveOsci2.OnChange(func(v string) {
		var oscitype OsciType
		switch v {
		case "rect":
			oscitype = OsciRect
		default:
			oscitype = OsciSin
		}
		k.instr.SetOsciType(2, oscitype)
	})
	k.waveOsci2.SetPos(pixel.V(450, 280))
	k.waveOsci2.SetAlignment(nodes.AlignmentCenterLeft)
	k.AddChild(k.waveOsci2)

	mididropdown := NewMidiDeviceDropDown()
	mididropdown.SetPos(pixel.V(790, 590))
	mididropdown.SetAlignment(nodes.AlignmentTopRight)
	mididropdown.OnChange(func(v string, _ string) {
		devid, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		if devid > 0 {
			if k.portIn != nil {
				k.portIn.Close()
				k.portIn = nil
			}
			portIn, err := portmidi.NewInputStream(portmidi.DeviceID(devid), 256)
			if err != nil {
				fmt.Println("opening midi device failed")
				return
			}
			k.portIn = portIn
		}
	})
	k.AddChild(mididropdown)
}

func (k *keyboardExample) Mount() {
	GetOutput().SetReadable(k.readable)
	Start()
}

func (k *keyboardExample) Unmount() {
	if k.portIn != nil {
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

	if k.portIn != nil {
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

func (k *keyboardExample) Draw(win pixel.Target, mat pixel.Matrix) {
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
}

func init() {
	k := &keyboardExample{
		BaseNode: *nodes.NewBaseNode("keyboard"),
	}
	k.Self = k
	AddExample("Keyboard", k)
}
