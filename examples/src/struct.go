package examples

import (
	"go-audio-service/snd"
)

type Example struct {
	Id    int
	Name  string
	Scene INode
}

var exampleList []*Example
var counter int = 1
var output snd.IOutput

func AddExample(name string, scene INode) {
	exampleList = append(exampleList, &Example{
		Id:    counter,
		Name:  name,
		Scene: scene,
	})
	scene.Init()
	counter++
}

func GetExamples() []*Example {
	return exampleList
}

func GetOutput() snd.IOutput {
	return output
}

func RunExample(id int) {
	for _, example := range exampleList {
		if example.Id == id {
			SetRoot(example.Scene)
		}
	}
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
