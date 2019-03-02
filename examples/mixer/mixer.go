package main

import (
	"fmt"
	"go-audio-service/generators"
	"go-audio-service/mix"
	"go-audio-service/snd"
	"time"
)

func main() {
	out, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	m := mix.NewMixer(44000)
	m.SetOutput(out)
	gen1 := generators.NewRect(44000, 440)
	ch1 := mix.NewChannel(44000)
	ch1.SetMixer(m)
	ch1.SetGain(0.2)
	ch1.SetPan(.2)
	gen1.SetOutput(ch1)

	gen2 := generators.NewRect(44000, 600)
	ch2 := m.GetChannel()
	ch2.SetGain(0.1)
	ch2.SetPan(.2)
	gen2.SetOutput(ch2)

	gen3 := generators.NewRect(44000, 900)
	ch3 := mix.NewChannel(44000)
	ch3.SetMixer(m)
	ch3.SetGain(0.2)
	ch3.SetPan(.2)
	gen3.SetOutput(ch3)

	err = out.Start()
	if err != nil {
		panic(err)
	}
	gen1.Start()
	gen2.Start()
	gen3.Start()
	time.Sleep(2 * time.Second)
	fmt.Println(".")
	gen1.Stop()
	time.Sleep(time.Second)
	fmt.Println(".")
	gen2.Stop()
	time.Sleep(time.Second)
	fmt.Println(".")
	gen3.Stop()
	_ = out.Stop()
}
