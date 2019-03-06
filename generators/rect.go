package generators

import (
	"go-audio-service/snd"
)

type Generator interface {
	SetOutput(f snd.Input)
	Start()
	Stop()
}

type Rect struct {
	snd.BasicReadable
	samplerate uint32
	high       bool
	current    int
	max        int
	done       chan struct{}
	running    bool
}

func NewRect(samplerate uint32, freq int) *Rect {
	return &Rect{
		samplerate: samplerate,
		high:       false,
		current:    0,
		max:        int(samplerate) / freq,
		done:       make(chan struct{}),
		running:    false,
	}
}

func (r *Rect) Read(samples *snd.Samples) int {
	length := len(samples.Frames)
	var v float32
	for i := 0; i < length; i++ {
		if r.running {
			if r.high {
				v = 0.5
			} else {
				v = -0.5
			}
		} else {
			v = 0.0
		}
		r.current++
		if r.current >= r.max {
			r.high = !r.high
			r.current = 0
		}
		samples.Frames[i].L = v
		samples.Frames[i].R = v
	}
	return length
}

func (r *Rect) Start() {
	r.running = true
}

func (r *Rect) Stop() {
	r.running = false
}
