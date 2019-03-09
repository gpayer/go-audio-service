package notes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndRemoveNotes(t *testing.T) {
	assert := assert.New(t)

	m := NewNoteMultiplexer()

	m.SendNoteEvent(NewNoteEvent(Pressed, Note(C, 3), 0.2))
	m.SendNoteEvent(NewNoteEvent(Pressed, Note(G, 3), 0.2))
	assert.Len(m.activeNotes, 2)
	m.SendNoteEvent(NewNoteEvent(Released, Note(C, 3), 0.2))
	m.SendNoteEvent(NewNoteEvent(Released, Note(G, 3), 0.2))
	assert.Len(m.activeNotes, 0)

	m.SendNoteEvent(NewNoteEvent(Pressed, Note(C, 3), 0.1))
	m.SendNoteEvent(NewNoteEvent(Pressed, Note(E, 3), 0.1))
	m.SendNoteEvent(NewNoteEvent(Pressed, Note(G, 3), 0.1))
	assert.Len(m.activeNotes, 3)
	m.SendNoteEvent(NewNoteEvent(Pressed, Note(C, 2), 0.2))
	assert.Len(m.activeNotes, 4)
	m.SendNoteEvent(NewNoteEvent(Released, Note(C, 3), 0.0))
	m.SendNoteEvent(NewNoteEvent(Released, Note(E, 3), 0.0))
	m.SendNoteEvent(NewNoteEvent(Released, Note(G, 3), 0.0))
	assert.Len(m.activeNotes, 1)
	m.SendNoteEvent(NewNoteEvent(Released, Note(C, 2), 0.0))
	assert.Len(m.activeNotes, 0)
}
