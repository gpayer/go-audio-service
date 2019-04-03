package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"pixelext/nodes"

	"github.com/faiface/pixel"
)

type sinExample struct {
	nodes.BaseNode
	totalTime float64
	gain      snd.Readable
}

func (s *sinExample) Init() {
	txt := nodes.NewText("txt", "basic")
	txt.SetZeroAlignment(nodes.AlignmentTopLeft)
	txt.SetPos(pixel.V(20, 580))
	txt.Printf("Sin example")
	s.AddChild(txt)

	sin := generators.NewSin(880)

	fminput, _ := sin.GetInput("fm")
	fmmod := generators.NewSin(880)
	fmgain := filters.NewGain(.9)
	fmgain.SetReadable(fmmod)
	fminput.SetReadable(fmgain)

	aminput, _ := sin.GetInput("am")
	ammod := generators.NewSin(3)
	amgain := filters.NewGain(.2)
	amgain.SetReadable(ammod)
	aminput.SetReadable(amgain)

	cont := notes.NewContinuousNote(notes.Note(notes.C, 3))
	cont.SetReadable(sin)

	gain := filters.NewGain(.3)
	gain.SetReadable(cont)
	s.gain = gain
}

func (s *sinExample) Mount() {
	s.totalTime = 0
	GetOutput().SetReadable(s.gain)
	Start()
}

func (s *sinExample) Unmount() {
	Stop()
}

func (s *sinExample) Update(dt float64) {
	s.totalTime += dt
	if s.totalTime >= 1.0 {
		SwitchScene("main")
	}
}

func init() {
	s := &sinExample{
		BaseNode:  *nodes.NewBaseNode("sin"),
		totalTime: 0,
	}
	s.Self = s
	AddExample("Sin", s)
}
