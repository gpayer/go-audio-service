package main

import (
	"fmt"
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

	cutoffValue.Start()
	rect.Start()
	err = output.Start()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		cutoff -= 7.0
		cutoffValue.Value = cutoff
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(500 * time.Millisecond)

	cutoffValue.Stop()
	rect.Stop()
	_ = output.Stop()
}
