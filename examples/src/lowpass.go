package examples

import (
	"fmt"
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/snd"
	"time"

	"github.com/faiface/pixel/pixelgl"
)

func runLowpass(output snd.IOutput, _ *pixelgl.Window) error {
	var cutoff float32 = 800.0

	lowpass := filters.NewLowPass(44000, cutoff, 1.0)
	lowpass.SetOutput(output)
	rect := generators.NewRect(44000, 800)
	lowpass.SetReadable(rect)
	cutoffInput, ok := lowpass.GetInput("cutoff")
	if !ok {
		panic(fmt.Errorf("no cutoff input"))
	}
	cutoffValue := generators.NewConstant(44000, cutoff)
	cutoffValue.SetOutput(cutoffInput)

	err := output.Start()
	if err != nil {
		return err
	}
	for i := 0; i < 100; i++ {
		cutoff -= 7.0
		cutoffValue.Value = cutoff
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(500 * time.Millisecond)

	return output.Stop()
}

func init() {
	AddExample("Lowpass", runLowpass)
}
