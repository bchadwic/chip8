package display

import (
	"log"

	"github.com/bchadwic/chip8/emulator/display/emit"
	"github.com/gonutz/prototype/draw"
)

type Display interface {
	Clear()
	Start()
	Set(e emit.Emit, row, col uint8)
	Get(row, col uint8) emit.Emit
}

const (
	SCALE = 10
)

type display struct {
	rows, cols int
	screen     []emit.Emit
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

func (d *display) Start() {
	err := draw.RunWindow("CHIP-8", d.cols*SCALE, d.rows*SCALE, d.update)
	if err != nil {
		log.Fatal(err.Error())
	}
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
	i := col + uint8(d.cols)*row
	return d.screen[i]
}

func (d *display) update(window draw.Window) {
	for i := 0; i < d.rows*d.cols; i++ {
		pixel := d.screen[i]
		var c draw.Color
		if pixel == emit.ON {
			c = draw.White
		} else {
			c = draw.Black
		}
		row := i / d.cols
		col := i % d.cols
		window.FillRect(col*SCALE, row*SCALE, 10, 10, c)
	}
}
