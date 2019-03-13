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

func (p *BasicWritableProvider) AddInput(name string) *BasicConnector {
	c := &BasicConnector{}
	p.inputs[name] = c
	return c
}

type BasicConnector struct {
	r Readable
}

func (c *BasicConnector) SetReadable(r Readable) {
	c.r = r
}

func (c *BasicConnector) Read(samples *Samples) {
	if c.r != nil {
		c.r.Read(samples)
		samples.Valid = false
	} else {
		samples.Valid = false
	}
}

func (c *BasicConnector) ReadStateless(samples *Samples, freq float32, state *NoteState) {
	c.r.ReadStateless(samples, freq, state)
}
