package examples

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type valueEntry struct {
	desc string
	val  int
	txt  *text.Text
}

type ValueToggle struct {
	w, h     float64
	current  int
	values   []*valueEntry
	onchange func(val int)
}

func NewValueToggle(w, h float64, onchange func(val int)) *ValueToggle {
	return &ValueToggle{
		w: w, h: h,
		onchange: onchange,
		current:  0,
		values:   make([]*valueEntry, 0),
	}
}

func (vt *ValueToggle) AddValue(desc string, val int) {
	txt := text.New(pixel.ZV, FontService.Get("basic"))
	fmt.Fprintf(txt, desc)
	vt.values = append(vt.values, &valueEntry{
		desc: desc,
		val:  val,
		txt:  txt,
	})
}

func (vt *ValueToggle) Init()      {}
func (vt *ValueToggle) Mounted()   {}
func (vt *ValueToggle) Unmounted() {}
func (vt *ValueToggle) Update(win *pixelgl.Window, dt float32, mat pixel.Matrix) {
	if len(vt.values) > 0 {
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			pos := win.MousePosition()
			if pixel.R(0, 0, vt.w, vt.h).Contains(mat.Unproject(pos)) {
				vt.current++
				if vt.current >= len(vt.values) {
					vt.current = 0
				}
			}
			vt.onchange(vt.values[vt.current].val)
		}
		entry := vt.values[vt.current]
		entry.txt.Draw(win, mat)
	}
}

func (vt *ValueToggle) Current() int {
	return vt.current
}

func (vt *ValueToggle) SetCurrent(c int) {
	vt.current = c
}
