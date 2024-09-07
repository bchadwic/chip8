package main

import "fmt"

func main() {
	var x uint16 = 0x80
	fmt.Printf("%08b\n", x)
	fmt.Printf("%x\n", x)
	fmt.Printf("%d\n", x)

}
