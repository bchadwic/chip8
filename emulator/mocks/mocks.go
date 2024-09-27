package mocks

import (
	"github.com/bchadwic/chip8/emulator/display"
	"github.com/bchadwic/chip8/emulator/display/emit"
)

type TestDisplay struct {
	// inputs
	In_Clear bool

	In_GetRow, In_GetCol uint8

	In_SetEmit           emit.Emit
	In_SetRow, In_SetCol uint8

	// outputs
	Out_GetEmit                            emit.Emit
	Out_PixelsPixels                       []display.Pixel
	Out_WindowSizeInt1, Out_WindowSizeInt2 int
}

func (td *TestDisplay) Clear() {
	td.In_Clear = true
}

func (td *TestDisplay) Get(row, col uint8) emit.Emit {
	td.In_GetRow = row
	td.In_GetCol = col
	return td.Out_GetEmit
}

func (td *TestDisplay) Set(e emit.Emit, row, col uint8) {
	td.In_SetEmit = e
	td.In_SetRow = row
	td.In_SetCol = col
}

func (td *TestDisplay) Pixels() []display.Pixel {
	return td.Out_PixelsPixels
}

func (td *TestDisplay) WindowSize() (int, int) {
	return td.Out_WindowSizeInt1, td.Out_WindowSizeInt2
}
