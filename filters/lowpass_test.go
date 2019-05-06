package filters

import (
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/snd"
	"testing"
)

func BenchmarkLowpassRead(b *testing.B) {
	b.ReportAllocs()
	cutoff := generators.NewConstant(44000, 100.0)
	rect := generators.NewRect(44000, 880.0)
	lowpass := NewLowPass(44000, 1.0, 1.0)
	lowpass.SetReadable(rect)
	cutoffW, _ := lowpass.GetInput("cutoff")
	cutoffW.SetReadable(cutoff)

	samples := snd.NewSamples(44000, 288)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lowpass.Read(samples)
	}
}
