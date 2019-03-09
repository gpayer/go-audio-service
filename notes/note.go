package notes

import "math"

type NoteValue float32

var BaseA4 NoteValue = 440.0

const (
	C      = "C"
	Csharp = "C#"
	Db     = "Db"
	D      = "D"
	Dsharp = "D#"
	Eb     = "Eb"
	E      = "E"
	F      = "F"
	Fsharp = "F#"
	Gb     = "Gb"
	G      = "G"
	Gsharp = "G#"
	Ab     = "Ab"
	A      = "A"
	Asharp = "A#"
	Hb     = "Hb"
	H      = "H"
	Bb     = "Bb"
	B      = "B"
)

var noteHalfSteps map[string]int = map[string]int{
	"C":  -9,
	"C#": -8,
	"Db": -8,
	"D":  -7,
	"D#": -6,
	"Eb": -6,
	"E":  -5,
	"F":  -4,
	"F#": -3,
	"Gb": -3,
	"G":  -2,
	"G#": -1,
	"Ab": -1,
	"A":  0,
	"A#": 1,
	"Hb": 1,
	"H":  2,
	"Bb": 1,
	"B":  2,
}

var twelthrootof2 float64 = 1.059463094359

func Note(name string, octave int) NoteValue {
	steps := 12*(octave-4) + noteHalfSteps[name]
	return BaseA4 * NoteValue(math.Pow(twelthrootof2, float64(steps)))
}
