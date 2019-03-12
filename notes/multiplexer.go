package notes

import (
	"go-audio-service/snd"
	"sync"
)

type noteInfo struct {
	timecode        uint32
	releaseTimecode uint32
	volume          float32
	on              bool
}

type NoteMultiplexer struct {
	mtx         sync.Mutex
	activeNotes map[NoteValue]*noteInfo
	tmp         *snd.Samples
	readable    snd.Readable
}

func NewNoteMultiplexer() *NoteMultiplexer {
	return &NoteMultiplexer{
		activeNotes: make(map[NoteValue]*noteInfo),
		tmp:         snd.NewSamples(44000, 128),
	}
}

func (n *NoteMultiplexer) SendNoteEvent(ev *NoteEvent) {
	n.mtx.Lock()
	defer n.mtx.Unlock()
	if ev.eventtype == Pressed {
		n.activeNotes[ev.note] = &noteInfo{timecode: 0, volume: ev.volume, on: true}
	} else {
		info, ok := n.activeNotes[ev.note]
		if ok {
			info.on = false
			info.releaseTimecode = info.timecode
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

	for note, info := range n.activeNotes {
		if !info.on && isNoteAware {
			noteaware.SetNoteReleased(info.releaseTimecode)
		}
		n.readable.ReadStateless(n.tmp, float32(note), info.timecode, info.on)
		info.timecode += uint32(length)

		for i := 0; i < length; i++ {
			samples.Frames[i].L += n.tmp.Frames[i].L * info.volume
			samples.Frames[i].R += n.tmp.Frames[i].R * info.volume
		}

		if isNoteAware {
			if noteaware.NoteEnded() {
				delete(n.activeNotes, note)
			}
		} else if !info.on {
			delete(n.activeNotes, note)
		}
	}
}

func (n *NoteMultiplexer) ReadStateless(samples *snd.Samples, freq float32, timecode uint32, _ bool) {
	n.Read(samples)
}
