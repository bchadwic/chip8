package drivers

import (
	"log"
	"strings"
	"sync"
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

	// display settings
	displayInitialized bool
	frameRate          int
	fill               bool
	color              draw.Color

	frame int
}

func Create(speaker speaker.Speaker, keypad keypad.Keypad, display display.Display) *driverContext {
	return &driverContext{
		speaker: speaker,
		keypad:  keypad,
		display: display,
		fill:    true,
		color:   draw.White,
	}
}

func (dc *driverContext) DisplaySettings(frameRate int, fill bool, color string) *driverContext {
	dc.displayInitialized = true
	dc.frameRate = frameRate
	dc.fill = fill
	switch strings.ToLower(color) {
	case "red":
		dc.color = draw.Red
	case "green":
		dc.color = draw.Green
	case "blue":
		dc.color = draw.Blue
	case "gray", "grey":
		dc.color = draw.Gray
	default:
		dc.color = draw.White
	}
	return dc
}

func (dc *driverContext) Start() {
	if !dc.displayInitialized {
		log.Fatal("display driver was not initialized")
	}
	rows, cols := dc.display.WindowSize()
	err := draw.RunWindow("CHIP-8", cols*display.SCALE, rows*display.SCALE, dc.update)
	if err != nil {
		log.Fatalf("an error occurred starting driver: %v", err)
	}
}

func (dc *driverContext) update(devices draw.Window) {
	// rate limit the updates
	time.Sleep(1 * time.Millisecond)
	dc.frame++
	var wg sync.WaitGroup
	wg.Add(3)

	go dc.renderDisplay(&wg, devices)
	go dc.readKeyboard(&wg, devices)
	go dc.playSpeakers(&wg, devices)

	wg.Wait()
}

func (dc *driverContext) renderDisplay(wg *sync.WaitGroup, window draw.Window) {
	defer wg.Done()
	for _, pixel := range dc.display.Pixels() {
		c := draw.Black
		if pixel.Status == emit.ON {
			c = dc.color
		}
		if dc.fill {
			window.FillRect(pixel.Col*display.SCALE, pixel.Row*display.SCALE, 10, 10, c)
		} else {
			window.DrawRect(pixel.Col*display.SCALE, pixel.Row*display.SCALE, 10, 10, c)
		}
	}
}

func (dc *driverContext) readKeyboard(wg *sync.WaitGroup, keyboard draw.Window) {
	defer wg.Done()
	chs := keyboard.Characters()
	for _, c := range chs {
		dc.keypad.Set(uint8(c))
	}
}

func (dc *driverContext) playSpeakers(wg *sync.WaitGroup, speakers draw.Window) {
	defer wg.Done()
	if dc.frame%dc.frameRate == 0 && dc.speaker.IsActive() {
		speakers.PlaySoundFile("beep.wav")
	}
}
