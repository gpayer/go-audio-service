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

type mp3Example struct {
	nodes.BaseNode
	multi *notes.NoteMultiplexer
	gain  snd.Readable
}

func (w *mp3Example) Init() {
	txt := ui.NewText("txt", "basic")
	txt.SetAlignment(nodes.AlignmentTopLeft)
	txt.SetPos(pixel.V(20, 580))
	txt.Printf("MP3 example")
	txt.Printf("\nPress space for sound\nPress Q to quit")
	w.AddChild(txt)

	samples, err := services.ResourceManager().LoadSample("samples/CYCdh_K4-Snr05.mp3")
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

func (w *mp3Example) Mount() {
	GetOutput().SetReadable(w.gain)
	Start()
}

func (w *mp3Example) Unmount() {
	Stop()
}

func (w *mp3Example) Update(dt float64) {
	if nodes.Events().JustPressed(pixelgl.KeyQ) {
		SwitchScene("main")
	}
	if nodes.Events().JustPressed(pixelgl.KeySpace) {
		w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Pressed, 880.0, .9))
	}
	if nodes.Events().JustReleased(pixelgl.KeySpace) {
		w.multi.SendNoteEvent(notes.NewNoteEvent(notes.Released, 880.0, .9))
	}
}

func init() {
	w := &mp3Example{
		BaseNode: *nodes.NewBaseNode("mp3"),
	}
	w.Self = w
	AddExample("MP3", w)
}
