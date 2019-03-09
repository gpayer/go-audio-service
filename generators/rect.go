package generators

import (
	"go-audio-service/snd"
)

type Generator interface {
	ReadStateless(samples *snd.Samples, freq float32, timecode uint32)
	SetGenerator(g Generator)
}

type Rect struct {
	snd.BasicReadable
	samplerate uint32
	high       bool
	current    int
	max        int
	running    bool
}

func NewRect(samplerate uint32, freq float32) *Rect {
	return &Rect{
		samplerate: samplerate,
		high:       false,
		current:    0,
		max:        int(float32(samplerate) / freq / 2),
		running:    false,
	}
}

func (r *Rect) Read(samples *snd.Samples) {
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
}

func (r *Rect) Start() {
	r.running = true
}

func (r *Rect) Stop() {
	r.running = false
}

func (r *Rect) ReadStateless(samples *snd.Samples, freq float32, timecode uint32, on bool) {
	length := len(samples.Frames)
	var v float32
	var max uint32
	if freq > 0 {
		max = uint32(float32(samples.SampleRate) / freq)
	} else {
		r.Read(samples)
		return
	}
	half := max / 2
	current := timecode % max
	for i := 0; i < length; i++ {
		if on {
			if current < half {
				v = 0.5
			} else {
				v = -0.5
			}
		} else {
			v = 0
		}
		samples.Frames[i].L = v
		samples.Frames[i].R = v
		current++
		if current >= max {
			current = 0
		}
	}
}
