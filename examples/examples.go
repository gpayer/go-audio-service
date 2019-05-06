package main

import (
	"flag"
	examples "github.com/gpayer/go-audio-service/examples/src"
	"log"
	"os"
	"github.com/gpayer/pixelext/nodes"
	"github.com/gpayer/pixelext/ui"
	"runtime"
	"runtime/pprof"

	"github.com/rakyll/portmidi"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type mainScene struct {
	nodes.BaseNode
	exampleListTxt *ui.Text
	numkeys        map[int]pixelgl.Button
}

func (m *mainScene) Init() {
	list := examples.GetExamples()
	m.exampleListTxt = ui.NewText("examplelist", "basic")
	m.exampleListTxt.Printf("press [ESC] to exit\n")
	for _, e := range list {
		m.exampleListTxt.Printf("%d: %s\n", e.Id, e.Name)
	}
	m.exampleListTxt.SetPos(pixel.V(20, 580))
	m.exampleListTxt.SetAlignment(nodes.AlignmentTopLeft)
	m.AddChild(m.exampleListTxt)

	m.numkeys = make(map[int]pixelgl.Button, 10)
	for i := 0; i <= 9; i++ {
		m.numkeys[i] = pixelgl.Button(int(pixelgl.Key0) + i)
	}
}

func (m *mainScene) Mounted() {
}

func (m *mainScene) Unmounted() {
}

func (m *mainScene) Update(dt float64) {
	for id, nk := range m.numkeys {
		if nodes.Events().JustPressed(nk) {
			examples.RunExample(id)
			nodes.SceneManager().Redraw()
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

	nodes.Events().SetWin(win)

	mainscene := &mainScene{
		BaseNode: *nodes.NewBaseNode("main"),
	}
	mainscene.Self = mainscene
	examples.AddScene("main", mainscene)
	nodes.SceneManager().SetRoot(mainscene)

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}

		nodes.SceneManager().Run(pixel.IM)
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
