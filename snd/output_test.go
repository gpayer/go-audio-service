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
