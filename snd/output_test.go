package snd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat2Bytes(t *testing.T) {
	assert := assert.New(t)
	bs := floatToBytes(0.0)
	assert.Equal(2, len(bs))
	assert.Zero(bs[0], bs[1])
	bs = floatToBytes(1.0)
	assert.Equal(byte(0xff), bs[0])
	assert.Equal(byte(0x7f), bs[1])
	bs = floatToBytes(-1.0)
	assert.Equal(byte(0x0), bs[0])
	assert.Equal(byte(0x80), bs[1])
	bs = floatToBytes(-1.0 / 32768)
	assert.Equal(byte(0xff), bs[0])
	assert.Equal(byte(0xff), bs[1])
}

func innerReadFn(o *Output, requestedSampleCount uint32, samples []byte) uint32 {
	// fmt.Println(requestedSampleCount, len(samples))
	input := Samples{
		SampleRate: o.samplerate,
		Frames:     make([]Sample, requestedSampleCount),
	}

	o.readable.Read(&input)

	offset := 0
	for i := uint32(0); i < requestedSampleCount; i++ {
		l := floatToBytes(input.Frames[i].L)
		samples[offset] = l[0]
		samples[offset+1] = l[1]
		r := floatToBytes(input.Frames[i].R)
		samples[offset+2] = r[0]
		samples[offset+3] = r[1]
		offset += 4
	}
	return requestedSampleCount
}

type simpleReadable struct {
	v float32
}

func (s *simpleReadable) Read(samples *Samples) {
	for i := 0; i < len(samples.Frames); i++ {
		samples.Frames[i].L = s.v
		samples.Frames[i].R = s.v
	}
}

func (s *simpleReadable) ReadStateless(samples *Samples, _ float32, _ uint32, _ bool) {
	s.Read(samples)
}

func BenchmarkInnerReadFn(b *testing.B) {
	o := &Output{}
	r := &simpleReadable{v: 23.0}
	o.readable = r

	samples := make([]byte, 223*4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = innerReadFn(o, 223, samples)
	}
}
