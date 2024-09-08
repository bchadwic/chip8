package main

import (
	"fmt"
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
	fmt.Println(fi.Size())
	rom := make([]uint8, fi.Size())
	_, err = f.Read(rom)
	if err != nil {
		log.Fatal(err.Error())
	}
	em := emulator.Create()
	em.Load(rom)
	em.Start()
}
