package display

import (
	"github.com/bchadwic/chip8/emulator/display/emit"
)

type Display interface {
	Clear()
	Set(e emit.Emit, row, col uint8)
	Get(row, col uint8) emit.Emit
	Pixels() chan Pixel
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

func (d *display) Set(e emit.Emit, row, col uint8) {
	i := int(row)*d.cols + int(col)
	d.screen[i] = e
}

func (d *display) Get(row, col uint8) emit.Emit {
	i := int(row)*d.cols + int(col)
	return d.screen[i]
}

func (d *display) Pixels() chan Pixel {
	pixels := make(chan Pixel, d.rows*d.cols)
	defer close(pixels)
	for i := 0; i < d.rows*d.cols; i++ {
		row := i / d.cols
		col := i % d.cols
		pixel := Pixel{Row: row, Col: col, Status: d.screen[i]}
		pixels <- pixel
	}
	return pixels
}

func (d *display) WindowSize() (int, int) {
	return d.rows, d.cols
}
