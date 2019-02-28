package main

import (
	"audiotest/snd"
	"time"
)

func main() {
	out, err := snd.NewOutput(44000)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	samples := &snd.Samples{SampleRate: 44000}
	samples.Frames = make([]snd.Sample, 500)
	for i := 0; i < 500; i++ {
		if i <= 250 {
			samples.Frames[i].L = -0.0
			samples.Frames[i].R = -0.2
		} else {
			samples.Frames[i].L = 0.0
			samples.Frames[i].R = 0.2
		}
	}

	for i := 0; i < 100; i++ {
		_ = out.Write(samples)
	}
	out.Start()
	time.Sleep(time.Second)
	out.Stop()
}
