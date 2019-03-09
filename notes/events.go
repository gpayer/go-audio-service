package notes

const (
	Pressed = iota
	Released
)

type NoteEvent struct {
	eventtype int
	note      NoteValue
	volume    float32
}

func NewNoteEvent(evtype int, note NoteValue, volume float32) *NoteEvent {
	return &NoteEvent{
		eventtype: evtype,
		note:      note,
		volume:    volume,
	}
}

type NoteEventReceiver interface {
	SendNoteEvent(ev *NoteEvent)
}
