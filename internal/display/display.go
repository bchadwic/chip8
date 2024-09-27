package display

import (
	"github.com/bchadwic/chip8/internal/display/emit"
)

type Display interface {
	Clear()
	Get(row, col uint8) emit.Emit
	Set(e emit.Emit, row, col uint8)
	Pixels() []Pixel
	WindowSize() (int, int)
}

const (
	SCALE = 10
)

type display struct {
	rows, cols int
	screen     []emit.Emit
}

type Pixel struct {
	Row, Col int
	Status   emit.Emit
}

func Create(rows, cols uint8) Display {
	irows, icols := int(rows), int(cols)
	display := &display{
		rows:   irows,
		cols:   icols,
		screen: make([]emit.Emit, irows*icols),
	}
	return display
}

func (d *display) Clear() {
	for i := 0; i < d.rows*d.cols; i++ {
		d.screen[i] = emit.OFF
	}
}

func (d *display) Get(row, col uint8) emit.Emit {
	i := int(row)*d.cols + int(col)
	return d.screen[i]
}

func (d *display) Set(e emit.Emit, row, col uint8) {
	i := int(row)*d.cols + int(col)
	d.screen[i] = e
}

func (d *display) Pixels() []Pixel {
	pixels := make([]Pixel, d.rows*d.cols)
	for i := 0; i < d.rows*d.cols; i++ {
		row := i / d.cols
		col := i % d.cols
		pixel := Pixel{Row: row, Col: col, Status: d.screen[i]}
		pixels[i] = pixel
	}
	return pixels
}

func (d *display) WindowSize() (int, int) {
	return d.rows, d.cols
}
