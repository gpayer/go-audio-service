package examples

import (
	"fmt"
	"go-audio-service/generators"
	"go-audio-service/mix"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

type noteShort struct {
	wait     time.Duration
	ch       int
	evtype   int
	notename string
	octave   int
	volume   float32
}

type adsrExample struct{}

func (a *adsrExample) Init() {
}

func (a *adsrExample) Mounted() {
	panic("not implemented")
}

func (a *adsrExample) Unmounted() {
	panic("not implemented")
}

func (a *adsrExample) Update(win *pixelgl.Window, dt float32) {
	panic("not implemented")
}

func createInstrument(instrtype int, a, d, s, r float32) *notes.NoteMultiplexer {
	var gen snd.Readable
	if instrtype == 1 {
		gen = generators.NewRect(44000, 440.0)
	} else if instrtype == 2 {
		sin := generators.NewSin(440)
		fm, _ := sin.GetInput("fm")
		fmmod := generators.NewSin(440)
		fm.SetReadable(fmmod)
		gen = sin
	}
	adsr1 := notes.NewAdsr(a, d, s, r)
	adsr1.SetReadable(gen)
	multi1 := notes.NewNoteMultiplexer()
	multi1.SetReadable(adsr1)

	return multi1
}

func runAdsr(output snd.IOutput, _ *pixelgl.Window) error {

	piece := []noteShort{
		{0, 1, notes.Pressed, "G", 2, 0.6},
		{0, 0, notes.Pressed, "C", 3, 1.0},
		{100, 0, notes.Released, "C", 3, 0.0},
		{0, 1, notes.Released, "G", 2, 0.0},
		{100, 0, notes.Pressed, "E", 3, 1.0},
		{250, 0, notes.Released, "E", 3, 0.0},
		{0, 0, notes.Pressed, "G", 3, 1.0},
		{0, 1, notes.Pressed, "C", 2, .6},
		{100, 1, notes.Released, "C", 2, 0.0},
		{500, 0, notes.Released, "G", 3, 0.0},
		{100, 1, notes.Pressed, "C", 4, 0.5},
		{100, 1, notes.Released, "C", 4, 0.0},
	}

	var instr []*notes.NoteMultiplexer
	instr = append(instr, createInstrument(1, 0.05, 0.3, 0.8, 0.5), createInstrument(2, 0.1, 2.0, 0.0, 0.0))

	mixer := mix.NewMixer(44000)
	ch1 := mixer.GetChannel()
	ch1.SetReadable(instr[0])
	ch1.SetGain(0.3)

	ch2 := mixer.GetChannel()
	ch2.SetReadable(instr[1])
	ch2.SetGain(0.4)

	mixer.SetGain(0.3)
	output.SetReadable(mixer)

	err := output.Start()
	if err != nil {
		return err
	}

	for _, n := range piece {
		if n.wait > 0 {
			time.Sleep(n.wait * time.Millisecond)
		}
		instr[n.ch].SendNoteEvent(notes.NewNoteEvent(n.evtype, notes.Note(n.notename, n.octave), n.volume))
		fmt.Printf("%d: %d %s %d\n", n.ch, n.evtype, n.notename, n.octave)
	}
	time.Sleep(1000 * time.Millisecond)

	return output.Stop()
}

func init() {
	AddExample("Adsr", &adsrExample{})
}
