package drivers

import (
	"log"
	"time"

	"github.com/bchadwic/chip8/emulator/display"
	"github.com/bchadwic/chip8/emulator/display/emit"
	"github.com/bchadwic/chip8/emulator/keypad"
	"github.com/bchadwic/chip8/emulator/speaker"
	"github.com/gonutz/prototype/draw"
)

type driverContext struct {
	speaker speaker.Speaker
	keypad  keypad.Keypad
	display display.Display
	frame   int
}

func Create(speaker speaker.Speaker, keypad keypad.Keypad, display display.Display) *driverContext {
	return &driverContext{
		speaker: speaker,
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
	time.Sleep(1 * time.Millisecond)
	driver.frame++
	w := make(chan string)
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
		w <- "display"
	}()

	go func() {
		chs := window.Characters()
		for _, c := range chs {
			driver.keypad.Set(uint8(c))
		}
		w <- "keypad"
	}()

	go func() {
		if driver.frame%4 == 0 && driver.speaker.IsActive() {
			window.PlaySoundFile("beep.wav")
		}
		w <- "sound"
	}()

	<-w
	<-w
	<-w
}
