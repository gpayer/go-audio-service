package examples

import (
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type FontServiceStruct struct {
	atlases map[string]*text.Atlas
}

func (f *FontServiceStruct) init() {
	f.atlases = make(map[string]*text.Atlas, 2)

	f.atlases["basic"] = text.NewAtlas(basicfont.Face7x13, text.ASCII)
}

func (f *FontServiceStruct) Get(name string) *text.Atlas {
	atlas, ok := f.atlases[name]
	if ok {
		return atlas
	}
	return f.atlases["basic"]
}

var FontService *FontServiceStruct

func init() {
	FontService = &FontServiceStruct{}
	FontService.init()
}
