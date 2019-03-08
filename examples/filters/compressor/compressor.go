package main

import (
	"go-audio-service/generators"
	"go-audio-service/mix"
	"go-audio-service/snd"
	"time"
)

func main() {
	output, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	var gain float32 = 0.1

	m := mix.NewMixer(44000)
	ch := m.GetChannel()
	ch.SetGain(gain)
	m.SetOutput(output)
	r := generators.NewRect(44000, 440)
	ch.SetReadable(r)

	r.Start()
	err = output.Start()
	if err != nil {
		panic(err)
	}

	var dg float32 = .9 / 100.0
	for i := 0; i < 100; i++ {
		time.Sleep(20 * time.Millisecond)
		gain += dg
		ch.SetGain(gain)
	}
	time.Sleep(time.Second)

	_ = output.Stop()
}
