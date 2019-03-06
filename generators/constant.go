package generators

import (
	"go-audio-service/snd"
)

type Constant struct {
	samplerate uint32
	Value      float32
	running    bool
}

func NewConstant(samplerate uint32, v float32) *Constant {
	return &Constant{
		samplerate: samplerate,
		Value:      v,
		running:    false,
	}
}

func (c *Constant) SetOutput(f snd.Writable) {
	f.SetReadable(c)
}

func (c *Constant) Start() {
	c.running = true
}

func (c *Constant) Stop() {
	c.running = false
}

func (c *Constant) Read(samples *snd.Samples) int {
	v := c.Value
	if !c.running {
		v = 0.0
	}
	for i := 0; i < len(samples.Frames); i++ {
		samples.Frames[i].L = v
		samples.Frames[i].R = v
	}
	return len(samples.Frames)
}
