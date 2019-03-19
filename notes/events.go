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

func (ev *NoteEvent) GetData() (int, NoteValue, float32) {
	return ev.eventtype, ev.note, ev.volume
}

type NoteEventReceiver interface {
	SendNoteEvent(ev *NoteEvent)
}
