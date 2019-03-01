package main

import (
	"go-audio-service/generators"
	"go-audio-service/snd"
	"time"
)

func main() {
	out, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	gen := generators.NewRect(44000, 440)
	gen.SetOutput(out)
	gen.Start()
	err = out.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	gen.Stop()
	_ = out.Stop()
}
