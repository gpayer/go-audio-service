package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"pixelext/nodes"
	"pixelext/ui"

	"github.com/faiface/pixel"
)

type rectExample struct {
	nodes.BaseNode
	totalTime float64
	gain      snd.Readable
}

func (r *rectExample) Init() {
	txt := ui.NewText("txt", "basic")
	txt.SetAlignment(nodes.AlignmentTopLeft)
	txt.SetPos(pixel.V(20, 580))
	txt.Printf("Rect example")
	r.AddChild(txt)

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

func (r *rectExample) Mount() {
	r.totalTime = 0
	GetOutput().SetReadable(r.gain)
	Start()
}

func (r *rectExample) Unmount() {
	Stop()
}

func (r *rectExample) Update(dt float64) {
	r.totalTime += dt
	if r.totalTime >= 1 {
		SwitchScene("main")
	}
}

func init() {
	r := &rectExample{
		BaseNode:  *nodes.NewBaseNode("rect"),
		totalTime: 0,
	}
	r.Self = r
	AddExample("Rect", r)
}
