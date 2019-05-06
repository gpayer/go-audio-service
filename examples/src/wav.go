package examples

import (
	"github.com/gpayer/go-audio-service/filters"
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/notes"
	"github.com/gpayer/go-audio-service/snd"
	"github.com/gpayer/pixelext/nodes"
	"github.com/gpayer/pixelext/services"
	"github.com/gpayer/pixelext/ui"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
)

type wavExample struct {
	nodes.BaseNode
	multi *notes.NoteMultiplexer
	gain  snd.Readable
}

func (w *wavExample) Init() {
	txt := ui.NewText("txt", "basic")
	txt.SetAlignment(nodes.AlignmentTopLeft)
	txt.SetPos(pixel.V(20, 580))
	txt.Printf("Wav example")
	txt.Printf("\nPress space for sound\nPress Q to quit")
	w.AddChild(txt)

	samples, err := services.ResourceManager().LoadSample("samples/CYCdh_K4-Snr05.wav")
	if err != nil {
		panic(err)
	}
	sampleplayer := generators.NewSample(samples)

	w.multi = notes.NewNoteMultiplexer()
	w.multi.SetReadable(sampleplayer)

	gain := filters.NewGain(.6)
	gain.SetReadable(w.multi)
	w.gain = gain
}

func (w *wavExample) Mount() {
	GetOutput().SetReadable(w.gain)
	Start()
}

func (w *wavExample) Unmount() {
	Stop()
}

func (w *wavExample) Update(dt float64) {
	if nodes.Events().JustPressed(pixelgl.KeyQ) {
		SwitchScene("main")
	}
	if nodes.Events().JustPressed(pixelgl.KeySpace) {
		w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, 880.0, .5))
	}
	if nodes.Events().JustReleased(pixelgl.KeySpace) {
		w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, 880.0, .5))
	}
}

func init() {
	w := &wavExample{
		BaseNode: *nodes.NewBaseNode("wav"),
	}
	w.Self = w
	AddExample("Wav", w)
}
