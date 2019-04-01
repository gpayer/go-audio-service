package examples

import (
	"pixelext/nodes"
	"github.com/faiface/pixel"
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/mix"

	"github.com/faiface/pixel/pixelgl"
)

type compressorExample struct {
	nodes.BaseNode
	totaltime float32
	comp      *filters.Compressor
	compstate *filters.CompressorState
	ch        *mix.Channel
}

func (c *compressorExample) Init() {
	m := mix.NewMixer(44000)
	c.ch = m.GetChannel()
	c.ch.SetGain(0.1)
	r := generators.NewRect(44000, 440)
	c.ch.SetReadable(r)
	compstate := filters.NewCompressorState()
	compstate.DefaultCompressor(44000)
	c.compstate = compstate
	comp := filters.NewCompressor(44000, compstate)
	comp.SetReadable(m)
	c.comp = comp
}

func (c *compressorExample) Mounted() {
	c.totaltime = 0
	c.ch.SetGain(0.1)
	c.compstate.DefaultCompressor(44000)
	output.SetReadable(c.comp)
	Start()
}

func (c *compressorExample) Unmounted() {
	Stop()
}

func (c *compressorExample) Update(win *pixelgl.Window, dt float32, mat pixel.Matrix) {
	c.totaltime += dt
	var gain float32 = .5 / 2.0 * c.totaltime
	c.ch.SetGain(gain)
	if c.totaltime > 2.0 {
		SwitchScene("main")
	}
}

func init() {
	AddExample("Compressor", &compressorExample{})
}
