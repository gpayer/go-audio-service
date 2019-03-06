package generators

import (
	"go-audio-service/snd"
)

type Constant struct {
	samplerate uint32
	Value      float32
	out        snd.Input
	running    bool
	done       chan struct{}
}

func NewConstant(samplerate uint32, v float32) *Constant {
	return &Constant{
		samplerate: samplerate,
		Value:      v,
		running:    false,
		done:       make(chan struct{}),
	}
}

func (c *Constant) SetOutput(f snd.Input) {
	c.out = f
}

func (c *Constant) Start() {
	if !c.running && c.out != nil {
		c.startGenerator()
	}
}

func (c *Constant) Stop() {
	if c.running {
		c.done <- struct{}{}
	}
}

func (c *Constant) startGenerator() {
	go func() {
		c.running = true
		for {
			select {
			case <-(c.done):
				c.running = false
				return
			default:
			}
			samples := &snd.Samples{SampleRate: c.samplerate}
			samples.Frames = make([]snd.Sample, 128)
			for i := 0; i < 128; i++ {
				samples.Frames[i].L = c.Value
				samples.Frames[i].R = c.Value
			}
			err := c.out.Write(samples)
			if err != nil {
				panic(err)
			}
		}
	}()
}
