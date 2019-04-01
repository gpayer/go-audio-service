package examples

import (
	"go-audio-service/generators"
	"go-audio-service/mix"
	"go-audio-service/notes"
	"go-audio-service/snd"
	"pixelext/nodes"

	"github.com/faiface/pixel"
)

type noteShort struct {
	wait     float32
	ch       int
	evtype   int
	notename string
	octave   int
	volume   float32
}

type adsrExample struct {
	nodes.BaseNode
	totaltime float32
	outrotime float32
	piece     []noteShort
	readable  snd.Readable
	instr     []*notes.NoteMultiplexer
	logtxt    *nodes.Text
}

func (a *adsrExample) initPiece() {
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
	eventtime := float32(0)

	for _, note := range piece {
		note.wait /= 1000.0
		eventtime += note.wait
		note.wait = eventtime
		a.piece = append(a.piece, note)
	}
}

func (a *adsrExample) Init() {
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

	a.readable = mixer
	a.instr = instr

	a.logtxt = nodes.NewText("", "basic")
	a.logtxt.Printf("ADSR Example")
	a.logtxt.SetPos(pixel.V(20, 580))
	a.logtxt.SetZeroAlignment(nodes.AlignmentTopLeft)
	a.AddChild(a.logtxt)
}

func (a *adsrExample) Mounted() {
	a.totaltime = 0
	a.outrotime = 0
	a.initPiece()
	GetOutput().SetReadable(a.readable)
	Start()
}

func (a *adsrExample) Unmounted() {
	Stop()
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

func (a *adsrExample) Update(dt float32) {
	a.totaltime += dt
	if len(a.piece) > 0 {
		n := a.piece[0]
		if a.totaltime >= n.wait {
			a.piece = a.piece[1:]
			a.instr[n.ch].SendNoteEvent(notes.NewNoteEvent(n.evtype, notes.Note(n.notename, n.octave), n.volume))
			a.logtxt.Printf("%d: %d %s %d\n", n.ch, n.evtype, n.notename, n.octave)
		}
	} else {
		a.outrotime += dt
		if a.outrotime > 1.0 {
			SwitchScene("main")
		}
	}
}

func init() {
	e := &adsrExample{
		BaseNode: *nodes.NewBaseNode("asdr"),
	}
	e.Self = e
	AddExample("Adsr", e)
}
