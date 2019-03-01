package generators

import (
	"go-audio-service/snd"
)

type Generator interface {
	SetOutput(f snd.Filter)
	Start()
	Stop()
}

type Rect struct {
	samplerate uint32
	out        snd.Filter
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

func (r *Rect) SetOutput(f snd.Filter) {
	r.out = f
}

func (r *Rect) startGenerator() {
	go func() {
		r.running = true
		for {
			select {
			case <-(r.done):
				r.running = false
				return
			default:
			}
			var v float32
			if r.high {
				v = 0.5
			} else {
				v = -0.5
			}
			samples := &snd.Samples{SampleRate: r.samplerate}
			for i := 0; r.current <= r.max && i < 512; i++ {
				samples.Add(snd.Sample{
					L: v,
					R: v,
				})
				r.current++
			}
			if r.current >= r.max {
				r.high = !r.high
				r.current = 0
			}
			err := r.out.Write(samples)
			if err != nil {
				panic(err)
			}
		}
	}()
}

func (r *Rect) Start() {
	if !r.running {
		r.startGenerator()
	}
}

func (r *Rect) Stop() {
	if r.running {
		r.done <- struct{}{}
	}
}
