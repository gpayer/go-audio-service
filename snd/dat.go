package snd

import (
	"fmt"
	"os"
	"time"
)

type DatWriter struct {
	samplerate uint32
	f          *os.File
	readable   Readable
	done       chan struct{}
}

func NewDatWriter(samplerate uint32, filepath string) (*DatWriter, error) {
	f, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	dat := &DatWriter{
		samplerate: samplerate,
		f:          f,
		done:       make(chan struct{}),
	}
	return dat, nil
}

func (dat *DatWriter) SetReadable(r Readable) {
	dat.readable = r
}

func (dat *DatWriter) Start() error {
	if dat.readable == nil {
		return fmt.Errorf("readable is nil")
	}
	go func() {
		samples := NewSamples(dat.samplerate, 128)
		timecode := 0
		ticker := time.NewTicker(time.Duration(1.0 / float64(dat.samplerate) * 1000000000.0 * 128.0))
		for true {
			select {
			case <-dat.done:
				ticker.Stop()
				return
			default:
			}
			dat.readable.Read(samples)
			for i := 0; i < 128; i++ {
				x := 1.0 / float32(dat.samplerate) * float32(timecode+i)
				y := samples.Frames[i].L
				_, err := dat.f.WriteString(fmt.Sprintf("%f %f\n", x, y))
				if err != nil {
					panic(err)
				}
			}
			timecode += 128
			<-ticker.C
		}
	}()
	return nil
}

func (dat *DatWriter) Stop() error {
	dat.done <- struct{}{}
	return nil
}

func (dat *DatWriter) Close() {
	_ = dat.f.Close()
}
