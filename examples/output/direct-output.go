package main

import (
	"go-audio-service/snd"
	"time"
)

func main() {
	out, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	samples := &snd.Samples{SampleRate: 44000}
	for i := 0; i < 500; i++ {
		var sample snd.Sample
		if i <= 250 {
			sample.L = -0.3
			sample.R = -0.3
		} else {
			sample.L = 0.3
			sample.R = 0.3
		}
		samples.Add(sample)
	}

	err = out.Start()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		_ = out.Write(samples)
	}
	time.Sleep(time.Second)
	_ = out.Stop()
}
