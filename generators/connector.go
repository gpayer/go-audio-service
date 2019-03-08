package generators

import "go-audio-service/snd"

type Connector struct {
	g Generator
}

func (c *Connector) ReadStateless(samples *snd.Samples, freq float32, timecode uint32) {
	c.g.ReadStateless(samples, freq, timecode)
}

func (c *Connector) SetGenerator(g Generator) {
	c.g = g
}

type GeneratorProvider interface {
	GetInput(name string) (Generator, bool)
}

type BasicGeneratorProvider struct {
	inputs map[string]Generator
}

func (p *BasicGeneratorProvider) InitBasicGeneratorProvider() {
	p.inputs = make(map[string]Generator)
}

func (p *BasicGeneratorProvider) GetInput(name string) (Generator, bool) {
	w, ok := p.inputs[name]
	return w, ok
}

func (p *BasicGeneratorProvider) AddInput(name string) *Connector {
	c := &Connector{}
	p.inputs[name] = c
	return c
}
