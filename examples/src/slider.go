package examples

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Slider struct {
	min, max, current float32
	dirty             bool
	w, h              float64
	canvas            *pixelgl.Canvas
	onchange          func(v float32)
}

func (s *Slider) Init() {
	s.dirty = true
}

func (s *Slider) Mounted() {
}

func (s *Slider) Unmounted() {
}

func (s *Slider) Update(win *pixelgl.Window, dt float32, mat pixel.Matrix) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		pos := mat.Unproject(win.MousePosition())
		if pixel.R(0, 0, s.w, s.h).Contains(pos) {
			s.current = s.min + (s.max-s.min)*float32(pos.X/s.w)
			s.dirty = true
			s.onchange(s.current)
		}
	}
	if s.canvas == nil {
		s.canvas = pixelgl.NewCanvas(pixel.R(0, 0, s.w, s.h))
		s.dirty = true
	}
	if s.dirty {
		s.canvas.Clear(colornames.Black)
		im := imdraw.New(nil)
		im.Color = colornames.White
		im.Push(pixel.V(0, 0), pixel.V(0, s.h), pixel.V(s.w, s.h), pixel.V(s.w, 0))
		im.Polygon(2)
		currentw := (s.w - 2) * float64(s.current/(s.max-s.min))
		im.Color = colornames.Skyblue
		im.Push(pixel.V(1, 1), pixel.V(currentw, 1), pixel.V(currentw, s.h-1), pixel.V(1, s.h-1))
		im.Polygon(0)
		im.Draw(s.canvas)
	}
	s.canvas.Draw(win, mat.Moved(pixel.V(s.w/2, s.h/2)))
}

func (s *Slider) SetValue(v float32) {
	if v >= s.min && v <= s.max {
		s.current = v
		s.dirty = true
	}
}

func (s *Slider) Value() float32 {
	return s.current
}

func (s *Slider) OnChange(fn func(v float32)) {
	s.onchange = fn
}

func NewSlider(w, h float64, min, max, current float32) *Slider {
	return &Slider{
		w: w, h: h,
		min: min, max: max, current: current,
		onchange: func(_ float32) {},
	}
}
