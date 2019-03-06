package snd

type Writable interface {
	SetReadable(r Readable)
}

type Readable interface {
	Read(samples *Samples) int
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
func (c *BasicConnector) Read(samples *Samples) int {
	return c.r.Read(samples)
}
