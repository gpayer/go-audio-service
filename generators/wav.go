package generators

import (
	"fmt"
	"go-audio-service/snd"
	"io"
	"os"

	"github.com/go-audio/wav"
)

type ReadSeekBuffer struct {
	Buf []byte
	Pos int64
}

func NewReadSeekBuffer(buf []byte) *ReadSeekBuffer {
	return &ReadSeekBuffer{
		Buf: buf,
	}
}

func (b *ReadSeekBuffer) Read(p []byte) (n int, err error) {
	panic("not implemented")
}

func (b *ReadSeekBuffer) Seek(offset int64, whence int) (int64, error) {
	var newpos int64
	switch whence {
	case io.SeekStart:
		newpos = offset
	case io.SeekCurrent:
		newpos += offset
	case io.SeekEnd:
		newpos = int64(len(b.Buf)) - offset
	}
	if newpos < 0 {
		return 0, fmt.Errorf("illegal file offset")
	}
	b.Pos = newpos
	return newpos, nil
}

type WavDecoder struct {
	decoder *wav.Decoder
	buf     *ReadSeekBuffer
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
	info, err := f.Stat()
	if err != nil {
		return err
	}
	buf := make([]byte, info.Size())
	_, err = f.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}

	w.buf = NewReadSeekBuffer(buf)
	w.decoder = wav.NewDecoder(w.buf)

	return nil
}

func (w *WavDecoder) Read(samples *snd.Samples) {
	panic("not implemented")
}

func (w *WavDecoder) ReadStateless(samples *snd.Samples, freq float32, state *snd.NoteState) {
	panic("not implemented")
}
