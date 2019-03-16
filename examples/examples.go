package main

import (
	"flag"
	"fmt"
	examples "go-audio-service/examples/src"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func run() {
	list := examples.GetExamples()

	cfg := pixelgl.WindowConfig{
		Title:  "go-audio-service examples",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	exampleRunning := false

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	exampleListTxt := text.New(pixel.ZV, basicAtlas)
	fmt.Fprintln(exampleListTxt, "press [ESC] to exit")
	for _, e := range list {
		fmt.Fprintf(exampleListTxt, "%d: %s\n", e.Id, e.Name)
	}

	var exampleRunningTxt *text.Text
	var done chan struct{}

	numkeys := make(map[int]pixelgl.Button, 10)
	for i := 0; i <= 9; i++ {
		numkeys[i] = pixelgl.Button(int(pixelgl.Key0) + i)
	}

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) {
			return
		}

		if exampleRunning {
			select {
			case <-done:
				exampleRunning = false
			default:
			}
		} else {
			for i, nk := range numkeys {
				if win.JustPressed(nk) {
					exampleRunning = true
					exampleRunningTxt = text.New(pixel.ZV, basicAtlas)
					fmt.Fprintf(exampleRunningTxt, "running example %d", i)
					for idx, e := range list {
						if e.Id == i {
							done = examples.RunExample(idx, win)
						}
					}
					break
				}
			}
		}
		win.Clear(colornames.Black)
		if exampleRunning {
			exampleRunningTxt.Draw(win, pixel.IM.Moved(pixel.V(20, win.Bounds().H()-20.0)))
		} else {
			exampleListTxt.Draw(win, pixel.IM.Moved(pixel.V(20, win.Bounds().H()-20.0)))
		}
		win.Update()
	}
}

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

	defer examples.Close()

	pixelgl.Run(run)
}
