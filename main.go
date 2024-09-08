package main

import (
	"log"
	"os"

	"github.com/bchadwic/chip8/emulator"
)

func main() {
	f, err := os.Open("ibm.ch8")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err.Error())
	}
	rom := make([]uint8, fi.Size())
	_, err = f.Read(rom)
	if err != nil {
		log.Fatal(err.Error())
	}
	em := emulator.Create()
	em.Load(rom)
	em.Start()
	// d := display.Create(32, 64)
	// d.Set(emit.ON, 3, 5)
	// d.Set(emit.ON, 0, 0)
	// d.Set(emit.ON, 25, 25)
	// d.Start()
}
