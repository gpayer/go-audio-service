package filters

import (
	"github.com/gpayer/go-audio-service/snd"
	"math"
)

type BiquadState struct {
	b0, b1, b2 float32
	a1, a2     float32
	xn1, xn2   snd.Sample
	yn1, yn2   snd.Sample
}

func (state *BiquadState) Process(input, output []snd.Sample) {
	if len(input) != len(output) {
		panic("input and output must have the same size")
	}
	size := len(input)

	b0 := state.b0
	b1 := state.b1
	b2 := state.b2
	a1 := state.a1
	a2 := state.a2
	xn1 := state.xn1
	xn2 := state.xn2
	yn1 := state.yn1
	yn2 := state.yn2

	// loop for each sample
	for n := 0; n < size; n++ {
		// get the current sample
		xn0 := input[n]

		// the formula is the same for each channel
		L :=
			b0*xn0.L +
				b1*xn1.L +
				b2*xn2.L -
				a1*yn1.L -
				a2*yn2.L
		R :=
			b0*xn0.R +
				b1*xn1.R +
				b2*xn2.R -
				a1*yn1.R -
				a2*yn2.R

		// save the result
		output[n] = snd.Sample{L: L, R: R}

		// slide everything down one sample
		xn2 = xn1
		xn1 = xn0
		yn2 = yn1
		yn1 = output[n]
	}

	// save the state for future processing
	state.xn1 = xn1
	state.xn2 = xn2
	state.yn1 = yn1
	state.yn2 = yn2
}

func (state *BiquadState) Reset() {
	state.xn1 = snd.Sample{L: 0, R: 0}
	state.xn2 = snd.Sample{L: 0, R: 0}
	state.yn1 = snd.Sample{L: 0, R: 0}
	state.yn2 = snd.Sample{L: 0, R: 0}
}

func (state *BiquadState) scale(amt float32) {
	state.b0 = amt
	state.b1 = 0.0
	state.b2 = 0.0
	state.a1 = 0.0
	state.a2 = 0.0
}

func (state *BiquadState) passThrough() {
	state.scale(1.0)
}

func (state *BiquadState) zero() {
	state.scale(0.0)
}

func (state *BiquadState) LowPass(rate uint32, cutoff, resonance float32) {
	nyquist := float32(rate) * 0.5
	cutoff /= nyquist

	if cutoff >= 1.0 {
		state.passThrough()
	} else if cutoff <= 0.0 {
		state.zero()
	} else {
		resonance = float32(math.Pow(10.0, float64(resonance)*0.05)) // convert resonance from dB to linear
		theta := math.Pi * 2.0 * float64(cutoff)
		alpha := float32(math.Sin(theta)) / (2.0 * resonance)
		cosw := float32(math.Cos(theta))
		beta := (1.0 - cosw) * 0.5
		a0inv := 1.0 / (1.0 + alpha)
		state.b0 = a0inv * beta
		state.b1 = a0inv * 2.0 * beta
		state.b2 = a0inv * beta
		state.a1 = a0inv * -2.0 * cosw
		state.a2 = a0inv * (1.0 - alpha)
	}
}
