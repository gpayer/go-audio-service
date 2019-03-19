package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/notes"
	"go-audio-service/snd"
)

type DoubleOsci struct {
	multi      *notes.NoteMultiplexer
	modgain    *filters.Gain
	adsr       *notes.Adsr
	osci       generators.Generator
	multimod   *notes.NoteMultiplexer
	oscimod    generators.Generator
	a, d, s, r float32
	modfactor  float32
	gainvalue  float32
}

func NewDoubleOsci(a, d, s, r, modfactor, modgain float32) *DoubleOsci {
	o := &DoubleOsci{
		a: a, d: d, s: s, r: r, modfactor: modfactor, gainvalue: modgain,
	}

	sin1 := generators.NewSin(440)
	o.osci = sin1
	fm, _ := sin1.GetInput("fm")
	fmmod := generators.NewSin(880)
	o.oscimod = fmmod
	o.modgain = filters.NewGain(modgain)
	o.modgain.SetReadable(fmmod)
	fmadsr := notes.NewAdsr(0.01, 0.01, 1.0, 0.01)
	fmadsr.SetReadable(o.modgain)
	o.multimod = notes.NewNoteMultiplexer()
	o.multimod.SetReadable(fmadsr)
	fm.SetReadable(o.multimod)
	o.adsr = notes.NewAdsr(a, d, s, r)
	o.adsr.SetReadable(o.osci)
	o.multi = notes.NewNoteMultiplexer()
	o.multi.SetReadable(o.adsr)
	return o
}

func (o *DoubleOsci) Read(samples *snd.Samples) {
	o.multi.Read(samples)
}

func (o *DoubleOsci) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	o.multi.Read(samples)
}

func (o *DoubleOsci) SendNoteEvent(ev *notes.NoteEvent) {
	o.multi.SendNoteEvent(ev)
	evtype, note, volume := ev.GetData()
	fmev := notes.NewNoteEvent(evtype, notes.NoteValue(float32(note)*o.modfactor), volume)
	o.multimod.SendNoteEvent(fmev)
}
