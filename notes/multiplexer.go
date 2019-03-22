package notes

import (
	"go-audio-service/snd"
	"sync"
)

type noteStruct struct {
	note  NoteValue
	state *snd.NoteState
}

type NoteMultiplexer struct {
	mtx         sync.Mutex
	activeNotes []noteStruct
	tmp         *snd.Samples
	readable    snd.Readable
}

func NewNoteMultiplexer() *NoteMultiplexer {
	return &NoteMultiplexer{
		activeNotes: make([]noteStruct, 0),
		tmp:         snd.NewSamples(44000, 128),
	}
}

func (n *NoteMultiplexer) SendNoteEvent(ev *NoteEvent) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	if ev.eventtype == Pressed {
		n.activeNotes = append(n.activeNotes, noteStruct{
			note:  ev.note,
			state: &snd.NoteState{Timecode: 0, Volume: ev.volume, On: true, Phase: 0},
		})
	} else {
		for _, nstruct := range n.activeNotes {
			if nstruct.note == ev.note && nstruct.state.On {
				nstruct.state.On = false
				nstruct.state.ReleaseTimecode = nstruct.state.Timecode
			}
		}
	}
}

func (n *NoteMultiplexer) SetReadable(r snd.Readable) {
	n.readable = r
}

func (n *NoteMultiplexer) Read(samples *snd.Samples) {
	n.mtx.Lock()
	defer n.mtx.Unlock()

	noteaware, isNoteAware := n.readable.(NoteAware)

	length := len(samples.Frames)
	if len(n.tmp.Frames) != length {
		n.tmp.Frames = make([]snd.Sample, length)
	}
	n.tmp.SampleRate = samples.SampleRate
	for i := 0; i < length; i++ {
		n.tmp.Frames[i].L = 0.0
		n.tmp.Frames[i].R = 0.0
	}

	var notesToDelete []int = make([]int, 0)
	for idx, nstruct := range n.activeNotes {
		note := nstruct.note
		info := nstruct.state
		n.readable.ReadStateless(n.tmp, float32(note), info)
		info.Timecode += uint32(length)

		for i := 0; i < length; i++ {
			samples.Frames[i].L += n.tmp.Frames[i].L * info.Volume
			samples.Frames[i].R += n.tmp.Frames[i].R * info.Volume
		}

		if isNoteAware {
			if noteaware.NoteEnded() {
				notesToDelete = append(notesToDelete, idx)
			}
		} else if !info.On {
			notesToDelete = append(notesToDelete, idx)
		}
	}
	for _, idx := range notesToDelete {
		n.activeNotes = delFromSlice(n.activeNotes, idx)
	}
}

func delFromSlice(a []noteStruct, i int) []noteStruct {
	if i < len(a)-1 {
		copy(a[i:], a[i+1:])
	}
	a[len(a)-1].state = nil
	a = a[:len(a)-1]
	return a
}

func (n *NoteMultiplexer) ReadStateless(samples *snd.Samples, freq float32, _ *snd.NoteState) {
	n.Read(samples)
}
