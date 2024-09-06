package emulator

import (
	"log"
	"time"
)

const (
	MEM_SIZE = 4096
	N1_MASK  = 0xF000
	N2_MASK  = 0x0F00
	N3_MASK  = 0x00F0
	N4_MASK  = 0x000F
)

type emulator struct {
	mem   []uint8
	pc    uint16
	i     uint16
	stack []uint16
}

func Create() *emulator {
	return &emulator{
		mem: make([]byte, MEM_SIZE),
	}
}

func (em *emulator) Start() {
	clock := time.NewTicker(3 * time.Microsecond)
	defer clock.Stop()

	for range clock.C {
		inst := em.fetch()
		em.decode(inst)
		em.execute()
	}
}

// fetch retrieves two bytes located at pc
// if two bytes are not available within
// the available memory, emulator will panic
func (em *emulator) fetch() uint16 {
	if em.pc > MEM_SIZE-2 {
		log.Fatal("pc out of memory bounds")
	}
	p1 := em.mem[em.pc]
	em.pc++
	p2 := em.mem[em.pc]
	em.pc++
	return (uint16(p1) << 8) | uint16(p2)
}

func (em *emulator) decode(inst uint16) {
	// n1 := inst & N1_MASK
	// n2 := inst & N2_MASK
	// n3 := inst & N3_MASK
	// n4 := inst & N4_MASK

}

func (em *emulator) execute() {
}
