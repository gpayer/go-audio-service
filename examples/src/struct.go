package examples

import (
	"go-audio-service/snd"
)

type Example struct {
	Id   int
	Name string
	fn   func(*snd.Output) error
}

var exampleList []*Example
var counter int = 1
var output *snd.Output

func AddExample(name string, fn func(*snd.Output) error) {
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

func RunExample(id int) {
	err := exampleList[id].fn(output)
	if err != nil {
		panic(err)
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

func init() {
	out, err := snd.NewOutput(44000, 512)
	if err != nil {
		panic(err)
	}
	output = out
}
