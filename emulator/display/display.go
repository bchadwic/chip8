package display

import (
	"fmt"
	"log"

	"github.com/bchadwic/chip8/emulator/display/emit"
	"github.com/gonutz/prototype/draw"
)

type Display interface {
	Clear()
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
	go func() {
		err := draw.RunWindow("CHIP-8", icols*SCALE, irows*SCALE, display.update)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	return display
}

func (d *display) Clear() {
	fmt.Println("size of screen!!")
	fmt.Println(len(d.screen))
	for i := 0; i < d.rows*d.cols; i++ {
		d.screen[i] = emit.OFF
	}
}

func (d *display) Set(e emit.Emit, row, col uint8) {
	i := col + uint8(d.cols)*row
	d.screen[i] = e
}

func (d *display) Get(row, col uint8) emit.Emit {
	i := col + uint8(d.cols)*row
	return d.screen[i]
}

func (d *display) update(window draw.Window) {
	for i := 0; i < d.rows*d.cols; i++ {
		pixel := d.screen[i]
		var fill draw.Color
		if pixel == emit.ON {
			fill = draw.White
		} else {
			fill = draw.Black
		}
		col := i % d.cols
		row := i / d.cols
		window.FillRect(col+(col*SCALE), row+(row*SCALE), 10, 10, fill)
	}
}
