package examples

import (
	"go-audio-service/snd"

	"github.com/faiface/pixel/pixelgl"
)

type Example struct {
	Id   int
	Name string
	fn   func(snd.IOutput, *pixelgl.Window) error
}

var exampleList []*Example
var counter int = 1
var output snd.IOutput

func AddExample(name string, fn func(snd.IOutput, *pixelgl.Window) error) {
	exampleList = append(exampleList, &Example{
		Id:   counter,
		Name: name,
		fn:   fn,
	})
	counter++
}

func GetExamples() []*Example {
	return exampleList
}

func RunExample(id int, win *pixelgl.Window) chan struct{} {
	done := make(chan struct{})
	go func() {
		err := exampleList[id].fn(output, win)
		if err != nil {
			panic(err)
		}
		done <- struct{}{}
	}()
	return done
}

func Start() {
	_ = output.Start()
}

func Stop() {
	_ = output.Stop()
}
func Close() {
	output.Close()
}

const (
	AudioOutput = iota
	DatOutput
)

func SetOutput(outputtype int, param string) {
	var out snd.IOutput
	var err error
	if outputtype == AudioOutput {
		out, err = snd.NewOutput(44000, 512)
	} else if outputtype == DatOutput {
		out, err = snd.NewDatWriter(44000, param)
	}
	if err != nil {
		panic(err)
	}
	output = out
}
