package examples

import (
	"github.com/faiface/pixel/pixelgl"
	"go-audio-service/generators"
	"go-audio-service/mix"
	"go-audio-service/snd"
	"time"
)

func runCompressor(output snd.IOutput, _ *pixelgl.Window) error {
	var gain float32 = 0.1

	m := mix.NewMixer(44000)
	ch := m.GetChannel()
	ch.SetGain(gain)
	m.SetOutput(output)
	r := generators.NewRect(44000, 440)
	ch.SetReadable(r)

	err := output.Start()
	if err != nil {
		return err
	}

	var dg float32 = .9 / 100.0
	for i := 0; i < 100; i++ {
		time.Sleep(20 * time.Millisecond)
		gain += dg
		ch.SetGain(gain)
	}
	time.Sleep(time.Second)

	return output.Stop()
}

func init() {
	AddExample("Compressor", runCompressor)
}
