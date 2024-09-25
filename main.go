package main

import (
	"flag"
	"log"
	"os"

	"github.com/bchadwic/chip8/emulator"
)

func main() {
	settings := &emulator.EmulatorSettings{}

	flag.IntVar(&settings.FrameRate, "r", 4, "frame refresh rate")
	flag.BoolVar(&settings.Fill, "l", false, "color fill pixels")
	flag.StringVar(&settings.Color, "c", "white", "color of pixels")
	flag.Parse()

	fname := flag.Arg(0)
	if fname == "" {
		log.Fatal("rom file not specified")
	}
	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("could not open rom file: %v", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatalf("could not determine rom size: %v", err)
	}
	rom := make([]uint8, fi.Size())
	_, err = f.Read(rom)
	if err != nil {
		log.Fatalf("could not read rom file: %v", err)
	}
	em := emulator.Create(settings)
	em.Load(rom)
	em.Start()
}
