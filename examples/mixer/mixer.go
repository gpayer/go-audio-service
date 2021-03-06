package main

import (
	"fmt"
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/mix"
	"github.com/gpayer/go-audio-service/snd"
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
	ch1 := m.GetChannel()
	ch1.SetGain(0.2)
	ch1.SetPan(-.5)
	ch1.SetReadable(gen1)

	gen2 := generators.NewRect(44000, 600)
	ch2 := m.GetChannel()
	ch2.SetGain(0.1)
	ch2.SetPan(.5)
	ch2.SetReadable(gen2)

	gen3 := generators.NewRect(44000, 220)
	ch3 := m.GetChannel()
	ch3.SetGain(0.2)
	ch3.SetPan(.5)
	ch3.SetReadable(gen3)

	gen1.Start()
	gen2.Start()
	gen3.Start()
	err = out.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
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
