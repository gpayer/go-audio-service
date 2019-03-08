package filters

import (
	"go-audio-service/generators"
	"go-audio-service/snd"
	"testing"
)

func BenchmarkCompressor(b *testing.B) {
	state := NewCompressorState()
	state.DefaultCompressor(44000)
	c := generators.NewConstant(44000, 0.95)
	comp := NewCompressor(44000, state)
	comp.SetReadable(c)
	samples := snd.NewSamples(44000, 223)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = comp.Read(samples)
	}
}
