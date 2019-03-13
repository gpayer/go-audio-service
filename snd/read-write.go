package snd

type Writable interface {
	SetReadable(r Readable)
}

type Readable interface {
	Read(samples *Samples)
	ReadStateless(samples *Samples, freq float32, state *NoteState)
}

type NoteState struct {
	Timecode        uint32
	ReleaseTimecode uint32
	Volume          float32
	On              bool
}

var EmptyNoteState = &NoteState{
	Timecode:        0,
	ReleaseTimecode: 0,
	Volume:          0,
	On:              true,
}

type BasicReadable struct{}

type WritableProvider interface {
	GetInput(name string) (Writable, bool)
}

type BasicWritableProvider struct {
	inputs map[string]Writable
}

func (p *BasicWritableProvider) InitBasicWritableProvider() {
	p.inputs = make(map[string]Writable)
}

func (p *BasicWritableProvider) GetInput(name string) (Writable, bool) {
	w, ok := p.inputs[name]
	return w, ok
}

func (p *BasicWritableProvider) AddInput(name string, defaultValue float32) *BasicConnector {
	c := &BasicConnector{
		defaultValue: defaultValue,
		samples:      NewSamples(44000, 128),
	}
	for i := 0; i < 128; i++ {
		c.samples.Frames[i].L = defaultValue
		c.samples.Frames[i].R = defaultValue
	}
	p.inputs[name] = c
	return c
}

type BasicConnector struct {
	r            Readable
	samples      *Samples
	defaultValue float32
}

func (c *BasicConnector) SetReadable(r Readable) {
	c.r = r
}

func (c *BasicConnector) Read(samples *Samples) {
	c.ReadStateless(samples, 0, EmptyNoteState)
}

func (c *BasicConnector) prepareBuffer(samplerate uint32, length int) {
	if c.samples.SampleRate != samplerate || len(c.samples.Frames) != length {
		c.samples = NewSamples(samplerate, length)
		for i := 0; i < length; i++ {
			c.samples.Frames[i].L = c.defaultValue
			c.samples.Frames[i].R = c.defaultValue
		}
	}
}

func (c *BasicConnector) ReadStateless(samples *Samples, freq float32, state *NoteState) {
	c.prepareBuffer(samples.SampleRate, len(samples.Frames))
	if c.r != nil {
		c.r.ReadStateless(samples, freq, state)
	}
}

func (c *BasicConnector) ReadBuffered(samplerate uint32, length int, freq float32, state *NoteState) *Samples {
	c.prepareBuffer(samplerate, length)
	c.ReadStateless(c.samples, freq, state)
	return c.samples
}
