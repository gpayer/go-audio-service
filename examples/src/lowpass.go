package examples

import (
	"pixelext/nodes"
	"github.com/faiface/pixel"
	"fmt"
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/snd"

	"github.com/faiface/pixel/pixelgl"
)

type lowpassExample struct {
	nodes.BaseNode
	totaltime   float32
	readable    snd.Readable
	cutoff      float32
	cutoffValue *generators.Constant
}

func (l *lowpassExample) Init() {
	l.cutoff = 800.0

	lowpass := filters.NewLowPass(44000, l.cutoff, 1.0)
	rect := generators.NewRect(44000, 800)
	lowpass.SetReadable(rect)
	l.readable = lowpass
	cutoffInput, ok := lowpass.GetInput("cutoff")
	if !ok {
		panic(fmt.Errorf("no cutoff input"))
	}
	cutoffValue := generators.NewConstant(44000, l.cutoff)
	cutoffValue.SetOutput(cutoffInput)
	l.cutoffValue = cutoffValue
}

func (l *lowpassExample) Mounted() {
	l.cutoff = 800.0
	l.totaltime = 0
	GetOutput().SetReadable(l.readable)
	Start()
}

func (l *lowpassExample) Unmounted() {
	Stop()
}

func (l *lowpassExample) Update(win *pixelgl.Window, dt float32, mat pixel.Matrix) {
	l.totaltime += dt
	if l.totaltime > 1.0 {
		SwitchScene("main")
	}
	if l.cutoff > 100.0 {
		l.cutoff -= 700.0 / 0.5 * dt
		l.cutoffValue.Value = l.cutoff
	}
}

func init() {
	AddExample("Lowpass", &lowpassExample{totaltime: 0})
}
