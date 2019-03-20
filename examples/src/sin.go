package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"

	"github.com/faiface/pixel"

	"github.com/faiface/pixel/pixelgl"
)

type sinExample struct {
	totalTime float32
	gain      snd.Readable
}

func (s *sinExample) Init() {
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

func (s *sinExample) Mounted() {
	s.totalTime = 0
	GetOutput().SetReadable(s.gain)
	Start()
}

func (s *sinExample) Unmounted() {
	Stop()
}

func (s *sinExample) Update(_ *pixelgl.Window, dt float32, mat pixel.Matrix) {
	s.totalTime += dt
	if s.totalTime >= 1.0 {
		SwitchScene("main")
	}
}

func init() {
	AddExample("Sin", &sinExample{totalTime: 0})
}
