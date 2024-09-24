package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/bchadwic/chip8/emulator"
)

func main() {
	fname := flag.String("f", "", "filename of rom")
	flag.Parse()
	if *fname == "" {
		log.Fatal(errors.New("rom file not specified"))
	}
	f, err := os.Open(*fname)
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
}
