package generators

import (
	"go-audio-service/snd"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstant(t *testing.T) {
	assert := assert.New(t)
	c := NewConstant(44000, 30.0)
	out := &snd.BufferedOutput{}

	c.SetOutput(out)
	c.Start()
	time.Sleep(1 * time.Millisecond)
	c.Stop()

	assert.True(len(out.Frames) >= 128)
	for _, fr := range out.Frames {
		assert.Equal(float32(30.0), fr.L)
		assert.Equal(float32(30.0), fr.R)
	}
}
