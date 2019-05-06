package main

import (
	"github.com/gpayer/go-audio-service/generators"
	"github.com/gpayer/go-audio-service/snd"
	"time"
)

func main() {
	out, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	gen := generators.NewRect(44000, 440)
	out.SetReadable(gen)
	err = out.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	_ = out.Stop()
}
