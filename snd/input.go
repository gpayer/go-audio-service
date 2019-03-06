package snd

import "fmt"

// Input defines the interface for filtering components
type Input interface {
	Write(samples *Samples) error
}

type InputProvider interface {
	GetInput(name string) (Input, bool)
}

type BufferedInputProvider struct {
	inputs map[string]*InputBuffer
}

type InputBuffer struct {
	channel    chan *Samples
	samplerate uint32
	buffer     []Sample
}

func (provider *BufferedInputProvider) Init() {
	provider.inputs = make(map[string]*InputBuffer)
}

func (provider *BufferedInputProvider) GetInput(name string) (*InputBuffer, bool) {
	input, ok := provider.inputs[name]
	return input, ok
}

func (provider *BufferedInputProvider) AddInput(name string, samplerate uint32) *InputBuffer {
	buf := NewInputBuffer(10, samplerate)
	provider.inputs[name] = buf
	return buf
}

func NewInputBuffer(size int, samplerate uint32) *InputBuffer {
	return &InputBuffer{
		channel:    make(chan *Samples, size),
		samplerate: samplerate,
	}
}

func (buffer *InputBuffer) Write(samples *Samples) error {
	if samples.SampleRate != buffer.samplerate {
		return fmt.Errorf("incompatible sample rates: %d != %d", samples.SampleRate, buffer.samplerate)
	}
	buffer.channel <- samples
	return nil
}

func (buffer *InputBuffer) Read(size int) *Samples {
	ok := true
	var samples *Samples
	for ok && len(buffer.buffer) < size {
		select {
		case samples = <-buffer.channel:
			ok = true
		default:
			ok = false
		}
		if ok {
			buffer.buffer = append(buffer.buffer, samples.Frames...)
		}
	}
	output := &Samples{SampleRate: buffer.samplerate}
	var effsize int
	if len(buffer.buffer) >= size {
		effsize = size
	} else {
		effsize = len(buffer.buffer)
	}
	output.Frames = buffer.buffer[:effsize]
	buffer.buffer = buffer.buffer[effsize:]
	return output
}
