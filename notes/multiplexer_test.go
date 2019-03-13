package notes

import (
	"go-audio-service/generators"
	"go-audio-service/snd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndRemoveNotes(t *testing.T) {
	assert := assert.New(t)

	m := NewNoteMultiplexer()
	r := generators.NewRect(44000, 880)
	m.SetReadable(r)
	samples := snd.NewSamples(44000, 256)

	m.SendNoteEvent(NewNoteEvent(Pressed, Note(C, 3), 0.2))
	m.SendNoteEvent(NewNoteEvent(Pressed, Note(G, 3), 0.2))
	m.Read(samples)
	assert.Len(m.activeNotes, 2)
	m.SendNoteEvent(NewNoteEvent(Released, Note(C, 3), 0.2))
	m.SendNoteEvent(NewNoteEvent(Released, Note(G, 3), 0.2))
	m.Read(samples)
	assert.Len(m.activeNotes, 0)

	m.SendNoteEvent(NewNoteEvent(Pressed, Note(C, 3), 0.1))
	m.SendNoteEvent(NewNoteEvent(Pressed, Note(E, 3), 0.1))
	m.SendNoteEvent(NewNoteEvent(Pressed, Note(G, 3), 0.1))
	m.Read(samples)
	assert.Len(m.activeNotes, 3)
	m.SendNoteEvent(NewNoteEvent(Pressed, Note(C, 2), 0.2))
	m.Read(samples)
	assert.Len(m.activeNotes, 4)
	m.SendNoteEvent(NewNoteEvent(Released, Note(C, 3), 0.0))
	m.SendNoteEvent(NewNoteEvent(Released, Note(E, 3), 0.0))
	m.SendNoteEvent(NewNoteEvent(Released, Note(G, 3), 0.0))
	m.Read(samples)
	assert.Len(m.activeNotes, 1)
	m.SendNoteEvent(NewNoteEvent(Released, Note(C, 2), 0.0))
	m.Read(samples)
	assert.Len(m.activeNotes, 0)
}
