package main

import (
	"flag"
	"fmt"
	examples "go-audio-service/examples/src"
)

func main() {
	outputtypeptr := flag.String("type", "audio", "output type: audio, dat")
	fileptr := flag.String("file", "", "output file for dat writer")

	flag.Parse()

	var outputtype int
	switch *outputtypeptr {
	case "audio":
		outputtype = examples.AudioOutput
	case "dat":
		outputtype = examples.DatOutput
	}

	examples.SetOutput(outputtype, *fileptr)

	list := examples.GetExamples()
	defer examples.Close()

	for true {
		fmt.Println(" 0: Quit")
		for _, example := range list {
			fmt.Printf("%2d: %s\n", example.Id, example.Name)
		}
		var idx int = -1
		fmt.Print("> ")
		fmt.Scanln(&idx)
		if idx == 0 {
			return
		} else {
			for i, example := range list {
				if idx == example.Id {
					examples.RunExample(i)
				}
			}
		}
	}
}
