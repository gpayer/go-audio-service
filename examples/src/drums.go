package examples

import (
	"go-audio-service/filters"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"pixelext/nodes"
	"pixelext/ui"

	"github.com/rakyll/portmidi"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
)

type drumsExample struct {
	nodes.BaseNode
	multi  *notes.NoteMultiplexer
	gain   snd.Readable
	portIn *portmidi.Stream
}

type mappingEntry struct {
	b pixelgl.Button
	n notes.NoteValue
	p string
}

var mapping []mappingEntry

func (w *drumsExample) Init() {
	txt := ui.NewText("txt", "basic")
	txt.SetAlignment(nodes.AlignmentTopLeft)
	txt.SetPos(pixel.V(20, 580))
	txt.Printf("Drums example")
	txt.Printf("\nPress keys for sound\nPress Q to quit")
	w.AddChild(txt)

	samplesmapper := NewNoteToSample()

	mapping = []mappingEntry{
		{pixelgl.KeyA, notes.Note(notes.C, 4), "samples/CYCdh_K4-Snr05.mp3"},
		{pixelgl.KeyS, notes.Note(notes.D, 4), "samples/CYCdh_K4-ClHat02.mp3"},
		{pixelgl.KeyD, notes.Note(notes.E, 4), "samples/CYCdh_K4-Kick02.mp3"},
		{pixelgl.KeyF, notes.Note(notes.F, 4), "samples/CYCdh_K4-Kick05.mp3"},
	}

	for _, m := range mapping {
		samplesmapper.AddSample(m.n, m.p)
	}

	w.multi = notes.NewNoteMultiplexer()
	w.multi.SetReadable(samplesmapper)

	gain := filters.NewGain(.6)
	gain.SetReadable(w.multi)
	w.gain = gain
}

func (w *drumsExample) Mount() {
	if *midiDeviceID > 0 {
		portIn, err := portmidi.NewInputStream(portmidi.DeviceID(*midiDeviceID), 256)
		if err != nil {
			panic(err)
		}
		w.portIn = portIn
	}
	GetOutput().SetReadable(w.gain)
	Start()
}

func (w *drumsExample) Unmount() {
	if *midiDeviceID > 0 && w.portIn != nil {
		w.portIn.Close()
		w.portIn = nil
	}
	Stop()
}

func (w *drumsExample) Update(dt float64) {
	if nodes.Events().JustPressed(pixelgl.KeyQ) {
		SwitchScene("main")
	}
	for _, m := range mapping {
		if nodes.Events().JustPressed(m.b) {
			w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, m.n, .9))
		}
		if nodes.Events().JustReleased(m.b) {
			w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, m.n, .9))
		}
	}

	if *midiDeviceID > 0 && w.portIn != nil {
		midiEvents, err := w.portIn.Read(128)
		if err == nil {
			for _, ev := range midiEvents {
				if ev.Status == 0x90 {
					w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, notes.MidiToNote(ev.Data1), float32(ev.Data2)/127.0))
				} else if ev.Status == 0x80 {
					w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, notes.MidiToNote(ev.Data1), 0.0))
				}
			}
		}
	}
}

func init() {
	w := &drumsExample{
		BaseNode: *nodes.NewBaseNode("drums"),
	}
	w.Self = w
	AddExample("Drums", w)
}
