package display

import (
	"testing"

	"github.com/bchadwic/chip8/emulator/display/emit"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	display := Create(3, 3)
	assert.NotNil(t, display)
}

func Test_Clear(t *testing.T) {
	display := &display{
		rows:   1,
		cols:   1,
		screen: []emit.Emit{emit.ON},
	}
	display.Clear()
	assert.Equal(t, emit.OFF, display.screen[0])
}

func Test_Get(t *testing.T) {
	display := &display{
		rows:   2,
		cols:   2,
		screen: []emit.Emit{emit.ON, emit.OFF, emit.ON, emit.OFF},
	}
	assert.Equal(t, emit.ON, display.Get(1, 0))
}

func Test_Set(t *testing.T) {
	display := &display{
		rows:   2,
		cols:   2,
		screen: []emit.Emit{emit.ON, emit.OFF, emit.ON, emit.OFF},
	}
	display.Set(emit.OFF, 1, 0)
	assert.Equal(t, emit.OFF, display.Get(1, 0))
}

func Test_Pixels(t *testing.T) {
	display := &display{
		rows:   2,
		cols:   2,
		screen: []emit.Emit{emit.ON, emit.OFF, emit.ON, emit.OFF},
	}
	pixels := display.Pixels()
	assert.Equal(t, len(display.screen), len(pixels))
}

func Test_WindowSize(t *testing.T) {
	display := &display{
		rows: 2,
		cols: 2,
	}
	rows, cols := display.WindowSize()
	assert.Equal(t, display.rows, rows)
	assert.Equal(t, display.cols, cols)
}
