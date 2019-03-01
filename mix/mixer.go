package mix

import (
	"go-audio-service/snd"
	"sync"
	"time"
)

type channelStruct struct {
	channel *Channel
	input   <-chan *snd.Samples
	buffer  []snd.Sample
}

// Mixer allows the mixing of different channels
type Mixer struct {
	mtx        sync.Mutex
	channels   []*channelStruct
	samplerate uint32
	gain       float32
	output     snd.Filter
	done       chan struct{}
	running    bool
}

// NewMixer creates a new Mixer instance
func NewMixer(samplerate uint32) *Mixer {
	return &Mixer{
		samplerate: samplerate,
		gain:       1.0,
		done:       make(chan struct{}),
		running:    false,
	}
}

func (m *Mixer) addChannel(ch *Channel) chan<- *snd.Samples {
	samplesCh := make(chan *snd.Samples, 512) // TODO: global configuration
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.channels = append(m.channels, &channelStruct{
		channel: ch,
		input:   samplesCh,
	})
	return samplesCh
}

// SetOutput sets the next filter in the output chain
func (m *Mixer) SetOutput(out snd.Filter) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.output = out
	if out != nil && !m.running {
		m.startWorker()
		m.running = true
	}
}

// SetGain sets the master gain value
func (m *Mixer) SetGain(gain float32) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.gain = gain
}

// Gain returns the master gain value
func (m *Mixer) Gain() float32 {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.gain
}

// Stop stops the mixer
func (m *Mixer) Stop() {
	if m.running {
		m.done <- struct{}{}
	}
}

func (m *Mixer) startWorker() {
	go func() {
		for {
			select {
			case <-m.done:
				m.mtx.Lock()
				defer m.mtx.Unlock()
				m.running = false
				return
			default:
			}

			m.mtx.Lock()
			minlen := 0
			for _, channel := range m.channels {
				select {
				case newSamples := <-(channel.input):
					channel.buffer = append(channel.buffer, newSamples.Frames...)
				default:
				}
				if minlen == 0 || len(channel.buffer) < minlen && len(channel.buffer) > 0 {
					minlen = len(channel.buffer)
				}
			}

			if minlen > 0 {
				buffer := make([]snd.Sample, minlen)
				for i := 0; i < minlen; i++ {
					buffer[i].L = 0.0
					buffer[i].R = 0.0
				}
				for _, channel := range m.channels {
					for i := 0; i < minlen; i++ {
						buffer[i].L += channel.buffer[i].L
						buffer[i].R += channel.buffer[i].R
					}
					channel.buffer = channel.buffer[minlen:]
				}
				err := m.output.Write(&snd.Samples{
					SampleRate: m.samplerate,
					Frames:     buffer,
				})
				if err != nil {
					m.mtx.Unlock()
					panic(err)
				}
				m.mtx.Unlock()
			} else {
				m.mtx.Unlock()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}
