package notes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameters(t *testing.T) {
	assert := assert.New(t)

	adsr := NewAdsr(0.5, 0.5, 0.5, 1.0)
	adsr.samplerate = 1000
	adsr.calcParameters()

	assert.Equal(500, int(adsr.t_decay))
	assert.Equal(1000, int(adsr.t_sustain))
	assert.Equal(float32(0.002), adsr.d_attack)
	assert.Equal(float32(-0.001), adsr.d_decay)
}

func TestReleaseParameters(t *testing.T) {
	assert := assert.New(t)

	adsr := NewAdsr(0.5, 0.5, 0.5, 1.0)
	adsr.samplerate = 1000
	adsr.calcParameters()

	adsr.calcRelease(1500) // in sustain
	assert.Equal(float32(0.5), adsr.releaseGain)
	assert.Equal(uint32(1500), adsr.releaseTimecode)
	assert.Equal(float32(-0.0005), adsr.d_release)

	adsr.calcRelease(750) // in decay
	assert.Equal(float32(0.75), adsr.releaseGain)
	assert.Equal(float32(-.75/1000.0), adsr.d_release)

	adsr.calcRelease(250) // in attack
	assert.Equal(float32(0.5), adsr.releaseGain)
	assert.Equal(float32(-.5/1000.0), adsr.d_release)
}
