package examples

import (
	"pixelext/nodes"
	"github.com/faiface/pixel"
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"

	"github.com/faiface/pixel/pixelgl"
)

type rectExample struct {
	nodes.BaseNode
	totalTime float32
	gain      snd.Readable
}

func (r *rectExample) Init() {
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
	r.gain = gain
}

func (r *rectExample) Mounted() {
	r.totalTime = 0
	GetOutput().SetReadable(r.gain)
	Start()
}

func (r *rectExample) Unmounted() {
	Stop()
}

func (r *rectExample) Update(win *pixelgl.Window, dt float32, mat pixel.Matrix) {
	r.totalTime += dt
	if r.totalTime >= 1 {
		SwitchScene("main")
	}
}

func init() {
	AddExample("Rect", &rectExample{totalTime: 0})
}
