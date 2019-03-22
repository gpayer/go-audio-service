package main

import (
	"flag"
	"fmt"
	examples "go-audio-service/examples/src"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/rakyll/portmidi"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type mainScene struct {
	basicAtlas     *text.Atlas
	exampleListTxt *text.Text
	numkeys        map[int]pixelgl.Button
}

func (m *mainScene) Init() {
	list := examples.GetExamples()
	m.basicAtlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
	m.exampleListTxt = text.New(pixel.ZV, m.basicAtlas)
	fmt.Fprintln(m.exampleListTxt, "press [ESC] to exit")
	for _, e := range list {
		fmt.Fprintf(m.exampleListTxt, "%d: %s\n", e.Id, e.Name)
	}

	m.numkeys = make(map[int]pixelgl.Button, 10)
	for i := 0; i <= 9; i++ {
		m.numkeys[i] = pixelgl.Button(int(pixelgl.Key0) + i)
	}
}

func (m *mainScene) Mounted() {
}

func (m *mainScene) Unmounted() {
}

func (m *mainScene) Update(win *pixelgl.Window, dt float32, mat pixel.Matrix) {
	m.exampleListTxt.Draw(win, pixel.IM.Moved(pixel.V(20, win.Bounds().H()-20.0)))

	for id, nk := range m.numkeys {
		if win.JustPressed(nk) {
			examples.RunExample(id)
			break
		}
	}
}

var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "go-audio-service examples",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	err = portmidi.Initialize()
	if err != nil {
		panic(err)
	}
	defer func() { _ = portmidi.Terminate() }()

	mainscene := &mainScene{}
	mainscene.Init()
	examples.AddScene("main", mainscene)
	examples.SetRoot(mainscene)

	last := time.Now()

	for !win.Closed() {
		dt := float32(time.Since(last).Seconds())
		last = time.Now()
		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}

		win.Clear(colornames.Black)
		examples.GetRoot().Update(win, dt, pixel.IM)
		win.Update()
	}

	if *memprofile != "" {
		fmemprofile, err := os.Create(*memprofile)
		if err != nil {
			panic(err)
		}
		defer fmemprofile.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(fmemprofile); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
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
