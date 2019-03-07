package filters

import "testing"

func BenchmarkBiquadLowpass(b *testing.B) {
	var cutoff float32 = 100.0
	state := &BiquadState{}
	for i := 0; i < b.N; i++ {
		state.LowPass(44000, cutoff, 1.0)
		cutoff += 7.0
		if cutoff > 10000.0 {
			cutoff = 100.0
		}
	}
}
