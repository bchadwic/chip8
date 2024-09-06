package main

import "fmt"

func main() {
	// em := emulator.Create()
	// em.Start()
	var i uint16 = 0b1001_0010_0000_0000
	// fmt.Printf("%x\n", i)
	// var j uint16 = 0xF
	// fmt.Printf("%x\n", j)
	// fmt.Printf("%x\n", i&j)
	var inversei uint16 = 0xF000
	fmt.Printf("%08b\n", i)
	fmt.Printf("%08b\n", inversei)
	fmt.Printf("%08b\n", i&inversei)
	fmt.Printf("%x\n", i&inversei)

}
