package examples

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/snd"
	"time"
)

func run(output *snd.Output) error {
	rect := generators.NewRect(44000, 440)

	fm, _ := rect.GetInput("fm")
	fmmod := generators.NewSin(200)
	fmgain := filters.NewGain(10)
	fmgain.SetReadable(fmmod)
	fm.SetReadable(fmgain)

	gain := filters.NewGain(0.3)
	gain.SetReadable(rect)
	output.SetReadable(gain)

	Start()
	time.Sleep(time.Second)
	Stop()

	return nil
}

func init() {
	AddExample("Rect", run)
}
