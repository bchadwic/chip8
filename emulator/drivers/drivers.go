package drivers

import (
	"fmt"
	"log"

	"github.com/bchadwic/chip8/emulator/display"
	"github.com/bchadwic/chip8/emulator/display/emit"
	"github.com/bchadwic/chip8/emulator/keypad"
	"github.com/gonutz/prototype/draw"
)

type driverContext struct {
	keypad  keypad.Keypad
	display display.Display
}

func Create(keypad keypad.Keypad, display display.Display) *driverContext {
	return &driverContext{
		keypad:  keypad,
		display: display,
	}
}

func (driver *driverContext) Start() {
	rows, cols := driver.display.WindowSize()
	err := draw.RunWindow(
		"CHIP-8",
		cols*display.SCALE,
		rows*display.SCALE,
		driver.update,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (driver *driverContext) update(window draw.Window) {
	w := make(chan any)
	defer close(w)

	go func() {
		// handle display
		for _, pixel := range driver.display.Pixels() {
			c := draw.Black
			if pixel.Status == emit.ON {
				c = draw.White
			}
			window.FillRect(
				pixel.Col*display.SCALE,
				pixel.Row*display.SCALE,
				10,
				10,
				c,
			)
		}
		w <- true
	}()

	driver.keypad.IsPressedFunc(func(kaddr uint8) bool {
		fmt.Println("kaddr", kaddr)
		fmt.Println("draw", draw.Key(kaddr))
		return window.IsKeyDown(draw.Key(kaddr))
	})
	driver.keypad.GetNextKeyFunc(func() uint8 {
		for {
			chs := window.Characters()
			if len(chs) == 0 {
				continue
			}
			return chs[0]
		}
	})
	<-w
}
