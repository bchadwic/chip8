package emulator

import (
	"math/rand"
	"testing"

	"github.com/bchadwic/chip8/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func testEmulator() *emulator {
	mem := make([]uint8, MEM_SIZE)
	// load fonts into memory
	for i := 0; i < len(fonts); i++ {
		mem[i+FONT_ADDR] = fonts[i]
	}

	return &emulator{
		registers: make([]uint8, REGISTERS),
		mem:       mem,
		stack:     make([]uint16, STACK_SIZE),
	}
}

func Test_Create(t *testing.T) {
	em := Create(&EmulatorSettings{})
	assert.NotNil(t, em)
	assert.Equal(t, len(em.registers), REGISTERS)
	assert.Equal(t, len(em.stack), STACK_SIZE)
}

func Test_Load(t *testing.T) {
	em := testEmulator()
	em.Load([]uint8{0xf, 0xf, 0xf})
	assert.Equal(t, em.pc, uint16(ROM_ADDR))
}

func Test_fetch(t *testing.T) {
	em := testEmulator()
	em.Load([]uint8{0x65, 0x05})
	fetchedInstruction, err := em.fetch()
	assert.Nil(t, err)
	assert.Equal(t, fetchedInstruction, uint16(0x6505))
}

func Test_cls(t *testing.T) {
	em := testEmulator()
	display := &mocks.TestDisplay{}
	em.display = display
	em.cls()
	assert.True(t, display.In_Clear)
}

func Test_ret(t *testing.T) {
	em := testEmulator()
	em.stack[2] = 0xFA
	em.sp = 3
	em.ret()
	assert.Equal(t, em.sp, uint8(2))
	assert.Equal(t, em.pc, uint16(0xFA))
}

func Test_jmp(t *testing.T) {
	em := testEmulator()
	em.pc = 0x233
	em.jmp(0x333)
	assert.Equal(t, em.pc, uint16(0x333))
}

func Test_call(t *testing.T) {
	em := testEmulator()
	em.sp = 3
	em.stack[3] = 0x222
	em.pc = 0x123

	em.call(0x333)
	assert.Equal(t, em.pc, uint16(0x333))
	assert.Equal(t, em.stack[3], uint16(0x123))
	assert.Equal(t, em.sp, uint8(4))
}

func Test_seqVxNN(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x32)
	em.pc = 3
	em.seqVxNN(3, 0x32)
	assert.Equal(t, em.pc, uint16(5))
	assert.Equal(t, em.registers[3], uint8(0x32))
}

func Test_sneVxNN(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x32)
	em.pc = 3
	em.sneVxNN(3, 0x32)
	assert.Equal(t, em.pc, uint16(3))
	assert.Equal(t, em.registers[3], uint8(0x32))
}

func Test_seqVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x32)
	em.registers[4] = uint8(0x32)
	em.pc = 3
	em.seqVxVy(3, 4)
	assert.Equal(t, em.pc, uint16(5))
	assert.Equal(t, em.registers[3], uint8(0x32))
}

func Test_ldVxKK(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x1)
	em.ldVxKK(3, 0x3)
	assert.Equal(t, em.registers[3], uint8(0x3))
}

func Test_addVxKK(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x1)
	em.addVxKK(3, 0x3)
	assert.Equal(t, em.registers[3], uint8(0x4))
}

func Test_ldVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x1)
	em.registers[4] = uint8(0x2)
	em.ldVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(0x2))
}

func Test_orVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x1)
	em.registers[4] = uint8(0x2)
	em.orVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(0x1)|uint8(0x2))
}

func Test_andVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x1)
	em.registers[4] = uint8(0x2)
	em.andVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(0x1)&uint8(0x2))
}

func Test_xorVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(0x1)
	em.registers[4] = uint8(0x2)
	em.xorVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(0x1)^uint8(0x2))
}

func Test_addVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(250)
	em.registers[4] = uint8(10)
	em.addVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(4))
	assert.Equal(t, em.registers[0xF], uint8(1))
}

func Test_subVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(250)
	em.registers[4] = uint8(10)
	em.subVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(0xf0))
	assert.Equal(t, em.registers[0xF], uint8(1))
}

func Test_shrVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(1)
	em.shrVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(0x0))
	assert.Equal(t, em.registers[0xF], uint8(1))
}

func Test_subnVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(3)
	em.registers[4] = uint8(5)
	em.shrVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(0x1))
	assert.Equal(t, em.registers[0xF], uint8(1))
}

func Test_shlVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(3)
	em.shlVxVy(3, 4)
	assert.Equal(t, em.registers[3], uint8(0x3)<<1)
	assert.Equal(t, em.registers[0xF], uint8(3)>>7)
}

func Test_sneVxVy(t *testing.T) {
	em := testEmulator()
	em.registers[3] = uint8(3)
	em.registers[4] = uint8(8)
	em.pc = 3
	em.sneVxVy(3, 4)
	assert.Equal(t, em.pc, uint16(5))
}

func Test_ldI(t *testing.T) {
	em := testEmulator()
	em.i = 3
	em.ldI(6)
	assert.Equal(t, em.i, uint16(6))
}

func Test_jmpV0(t *testing.T) {
	em := testEmulator()
	em.pc = 3
	em.registers[0] = 2
	em.jmpV0(6)
	assert.Equal(t, em.pc, uint16(8))
}

func Test_rndVxKK(t *testing.T) {
	em := testEmulator()
	em.pc = 3
	// deprecated, but I don't think I care yet
	rand.Seed(32)
	em.registers[3] = 2
	em.rndVxKK(3, 6)
	assert.Equal(t, em.pc, uint16(3))
}

func Test_ldVxDt(t *testing.T) {
	em := testEmulator()
	em.registers[3] = 2
	em.dt = 3
	em.ldVxDt(3)
	assert.Equal(t, em.registers[3], uint8(3))
}

func Test_ldDtVx(t *testing.T) {
	em := testEmulator()
	em.dt = 3
	em.registers[3] = 2
	em.ldDtVx(3)
	assert.Equal(t, em.dt, uint8(2))
}

func Test_ldStVx(t *testing.T) {
	em := testEmulator()
	em.st = 3
	em.registers[3] = 2
	em.ldStVx(3)
	assert.Equal(t, em.st, uint8(2))
}

func Test_addIVx(t *testing.T) {
	em := testEmulator()
	em.i = 3
	em.registers[3] = 2
	em.addIVx(3)
	assert.Equal(t, em.i, uint16(5))
	assert.Equal(t, em.registers[0xF], uint8(0))
}

func Test_ldFVx(t *testing.T) {
	em := testEmulator()
	em.i = 3
	em.registers[3] = 2
	em.ldFVx(3)
	assert.Equal(t, em.i, uint16(2)*5)
}

func Test_ldIVx(t *testing.T) {
	em := testEmulator()
	em.i = 3
	em.ldIVx(3)
	assert.Equal(t, em.i, uint16(4))
}

func Test_ldVxI(t *testing.T) {
	em := testEmulator()
	em.i = 3
	em.ldVxI(3)
	assert.Equal(t, em.i, uint16(4))
}
