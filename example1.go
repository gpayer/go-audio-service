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
	samples.Frames = make([]snd.Sample, 500)
	for i := 0; i < 500; i++ {
		if i <= 250 {
			samples.Frames[i].L = -0.3
			samples.Frames[i].R = -0.3
		} else {
			samples.Frames[i].L = 0.3
			samples.Frames[i].R = 0.3
		}
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
