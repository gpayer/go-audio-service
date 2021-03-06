package mix

import (
	"sync"

	"github.com/gpayer/go-audio-service/snd"
)

// Mixer allows the mixing of different channels
type Mixer struct {
	mtx        sync.Mutex
	channels   []*Channel
	samplerate uint32
	gain       float32
	running    bool
	tmp        *snd.Samples
}

// NewMixer creates a new Mixer instance
func NewMixer(samplerate uint32) *Mixer {
	return &Mixer{
		samplerate: samplerate,
		gain:       1.0,
		running:    true,
		tmp:        snd.NewSamples(samplerate, 256),
	}
}

func (m *Mixer) addChannel(ch *Channel) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.channels = append(m.channels, ch)
}

// SetOutput sets the next filter in the output chain
func (m *Mixer) SetOutput(out snd.Writable) {
	out.SetReadable(m)
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
	m.running = false
}

// GetChannel returns a new channel connected to this Mixer
func (m *Mixer) GetChannel() *Channel {
	ch := NewChannel(m.samplerate)
	m.addChannel(ch)
	return ch
}

// RemoveChannel removes an existing channel from the mixer
func (m *Mixer) RemoveChannel(ch *Channel) {
	for i, channel := range m.channels {
		if channel == ch {
			m.channels[i] = m.channels[len(m.channels)-1]
			m.channels[len(m.channels)-1] = nil
			m.channels = m.channels[:len(m.channels)-1]
			return
		}
	}
}

func (m *Mixer) Read(samples *snd.Samples) {
	m.ReadStateless(samples, 0, snd.EmptyNoteState)
}

func (m *Mixer) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	length := len(samples.Frames)
	if m.running {
		if len(samples.Frames) != len(m.tmp.Frames) {
			m.tmp = snd.NewSamples(m.samplerate, len(samples.Frames))
		}
		tmp := m.tmp
		for _, channel := range m.channels {
			if !channel.Enabled() {
				continue
			}
			channel.ReadStateless(tmp, freq, state)
			for i := 0; i < length; i++ {
				samples.Frames[i].L += tmp.Frames[i].L
				samples.Frames[i].R += tmp.Frames[i].R
			}
		}
	} else {
		for i := 0; i < length; i++ {
			samples.Frames[i].L = 0.0
			samples.Frames[i].R = 0.0
		}
	}
	for i := 0; i < length; i++ {
		samples.Frames[i].L *= m.gain
		samples.Frames[i].R *= m.gain
	}
}
