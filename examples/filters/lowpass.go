package main

import (
	"go-audio-service/filters"
	"go-audio-service/generators"
	"go-audio-service/snd"
	"time"
)

func main() {
	output, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	lowpass := filters.NewLowPass(44000, 800, 1.0)
	lowpass.SetOutput(output)
	rect := generators.NewRect(44000, 800)
	rect.SetOutput(lowpass)

	rect.Start()
	err = output.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	rect.Stop()
	_ = output.Stop()
}
