package generators

import (
	"go-audio-service/snd"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type WavDecoder struct {
	pos int
	buf *audio.Float32Buffer
}

func NewWavDecoder() *WavDecoder {
	w := &WavDecoder{}
	return w
}

func (w *WavDecoder) Load(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := wav.NewDecoder(f)
	intbuf, err := decoder.FullPCMBuffer()
	if err != nil {
		return err
	}
	w.buf = intbuf.AsFloat32Buffer()

	return nil
}

func (w *WavDecoder) Read(samples *snd.Samples) {
	state := snd.NoteState{
		On:       true,
		Timecode: uint32(w.pos),
		Volume:   0.5,
	}
	w.ReadStateless(samples, 0, &state)
	w.pos += w.buf.Format.NumChannels * len(samples.Frames)
}

func (w *WavDecoder) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	pos := w.buf.Format.NumChannels * int(state.Timecode)
	for i := 0; i < len(samples.Frames); i++ {
		if pos < len(w.buf.Data) {
			samples.Frames[i].L = w.buf.Data[pos]
			if w.buf.Format.NumChannels > 1 {
				samples.Frames[i].R = w.buf.Data[pos+1]
			}
			pos += w.buf.Format.NumChannels
		}
	}
}
