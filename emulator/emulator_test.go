package emulator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetch(t *testing.T) {
	tests := map[string]struct {
		pc        uint16
		p1, p2    uint8
		inst      uint16
		wantedErr error
	}{
		"t1": {
			pc:   20,
			p1:   0b1010_1010,
			p2:   0b1010_1010,
			inst: 0b1010_1010_1010_1010,
		},
		"t2": {
			pc:   500,
			p1:   0b1110_1110,
			p2:   0b0000_0011,
			inst: 0b1110_1110_0000_0011,
		},
		"t3": {
			pc:   600,
			p1:   0b1110_1110,
			p2:   0b0000_0011,
			inst: 0b1110_1110_0000_0011,
		},
		"t4": {
			pc:   MEM_SIZE - 2,
			p1:   0b1100_0011,
			p2:   0b0011_1100,
			inst: 0b1100_0011_0011_1100,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.pc = test.pc
			em.mem[test.pc] = test.p1
			em.mem[test.pc+1] = test.p2

			inst, err := em.fetch()
			assert.True(t, (err != nil) == (test.wantedErr != nil))
			if test.wantedErr != nil {
				assert.Equal(t, test.wantedErr, err)
				return
			}
			assert.Equal(t, test.inst, inst)
			assert.Equal(t, test.pc+2, em.pc)
		})
	}
}

func Test_cls(t *testing.T) {
	tests := map[string]struct {
		display []uint8
	}{
		"t1": {
			display: []uint8{
				0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
				1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
				0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
				1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
				0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
				1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
				0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0,
				1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1,
			},
		},
		"t2": {
			display: []uint8{
				1, 1, 6, 5, 0, 1, 0, 1, 0, 1, 1, 1, 8, 1, 0,
				1, 1, 1, 0, 1, 0, 1, 0, 8, 0, 1, 1, 1, 0, 1,
				1, 1, 2, 1, 0, 1, 0, 1, 1, 1, 0, 1, 9, 1, 0,
				1, 1, 1, 0, 1, 0, 3, 0, 5, 0, 1, 1, 0, 0, 1,
				1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0,
				1, 1, 1, 6, 1, 0, 9, 0, 1, 1, 1, 1, 3, 0, 1,
				1, 1, 1, 1, 0, 1, 0, 1, 7, 1, 1, 1, 2, 1, 0,
				1, 1, 1, 0, 1, 0, 1, 0, 7, 0, 1, 0, 3, 5, 1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.display = test.display

			em.cls()
			var zero uint8 = 0
			for _, e := range em.display {
				assert.Equal(t, zero, e)
			}
		})
	}
}

func Test_ret(t *testing.T) {
	tests := map[string]struct {
		pc        uint16
		sp        uint8
		stack     []uint16
		wantedErr error
	}{
		"t1": {
			pc: 0x03af,
			sp: 2,
			// 16 positions
			stack: []uint16{0x0dde, 0xdee, 0x0ddd, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		"t2": {
			pc: 0x03af,
			sp: 3,
			// 16 positions
			stack: []uint16{0x0dde, 0xdee, 0x0ddd, 0x0333, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		"t3": {
			pc: 0x0333,
			sp: 1,
			// 16 positions
			stack: []uint16{0x0dde, 0xdee, 0x0ddd, 0x0333, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		"t4": {
			pc: 0x0333,
			sp: 16,
			// 16 positions
			stack:     []uint16{0x0dde, 0xdee, 0x0ddd, 0x0333, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantedErr: fmt.Errorf("stack pointer out of bounds: %d", 16),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.pc = test.pc
			em.sp = test.sp
			em.stack = test.stack

			err := em.ret()
			assert.True(t, (err != nil) == (test.wantedErr != nil))
			if test.wantedErr != nil {
				assert.Equal(t, test.wantedErr, err)
				return
			}
			assert.Equal(t, em.pc, em.stack[em.sp+1])
		})
	}
}

func Test_jmp(t *testing.T) {
	tests := map[string]struct {
		addr uint16
	}{
		"t1": {
			addr: 0x0333,
		},
		"t2": {
			addr: 0x0334,
		},
		"t3": {
			addr: 0x0433,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.jmp(test.addr)
			assert.Equal(t, test.addr, em.pc)
		})
	}
}

func Test_call(t *testing.T) {
	tests := map[string]struct {
		pc        uint16
		sp        uint8
		stack     []uint16
		addr      uint16
		wantedErr error
	}{
		"t1": {
			pc:    0x0363,
			sp:    3,
			stack: []uint16{0x0dde, 0xdee, 0x0ddd, 0x0333, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		"t2": {
			pc:    0x0222,
			sp:    1,
			stack: []uint16{0x0dde, 0xdee, 0x0ddd, 0x0333, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		"t3": {
			pc:        0x0441,
			sp:        15,
			stack:     []uint16{0x0dde, 0xdee, 0x0ddd, 0x0333, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantedErr: errors.New("stack overflow"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.pc = test.pc
			em.sp = test.sp
			em.stack = test.stack
			err := em.call(test.addr)

			assert.True(t, (err != nil) == (test.wantedErr != nil))
			if test.wantedErr != nil {
				assert.Equal(t, test.wantedErr, err)
				return
			}
			assert.Equal(t, test.sp+1, em.sp)
			assert.Equal(t, test.pc, em.stack[test.sp+1])
			assert.Equal(t, test.addr, em.pc)
		})
	}
}

func Test_seqVxNN(t *testing.T) {
	tests := map[string]struct {
		pc        uint16
		registers []uint8
		x         uint16
		nn        uint16
		wantedInc bool
		wantedErr error
	}{
		"t1": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         2,
			nn:        0xff,
			// true because v[2] = 0xff, and nn = 0xff, so an increment of 2 should occur
			wantedInc: true,
		},
		"t2": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         2,
			nn:        0xfe,
			// false because v[2] = 0xff, and nn = 0xfe, so an increment should not occur
			wantedInc: false,
		},
		"t3": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         18,
			nn:        0xfe,
			wantedErr: errors.New("register index out of bounds"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.pc = test.pc
			em.registers = test.registers
			err := em.seqVxNN(test.x, test.nn)

			assert.True(t, (err != nil) == (test.wantedErr != nil))
			if test.wantedErr != nil {
				assert.Equal(t, test.wantedErr, err)
				return
			}
			if test.wantedInc {
				assert.Equal(t, test.pc+2, em.pc)
			} else {
				assert.Equal(t, test.pc, em.pc)
			}
		})
	}
}

func Test_sneVxNN(t *testing.T) {
	tests := map[string]struct {
		pc        uint16
		registers []uint8
		x         uint16
		nn        uint16
		wantedInc bool
		wantedErr error
	}{
		"t1": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         2,
			nn:        0xff,
			// false because v[2] = 0xff, and nn = 0xff, so an increment of 2 should occur
			wantedInc: false,
		},
		"t2": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         2,
			nn:        0xfe,
			// true because v[2] = 0xff, and nn = 0xfe, so an increment should not occur
			wantedInc: true,
		},
		"t3": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         18,
			nn:        0xfe,
			wantedErr: errors.New("register index out of bounds"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.pc = test.pc
			em.registers = test.registers
			err := em.sneVxNN(test.x, test.nn)

			assert.True(t, (err != nil) == (test.wantedErr != nil))
			if test.wantedErr != nil {
				assert.Equal(t, test.wantedErr, err)
				return
			}
			if test.wantedInc {
				assert.Equal(t, test.pc+2, em.pc)
			} else {
				assert.Equal(t, test.pc, em.pc)
			}
		})
	}
}

func Test_seqVxVy(t *testing.T) {
	tests := map[string]struct {
		pc        uint16
		registers []uint8
		x         uint16
		y         uint16
		wantedInc bool
		wantedErr error
	}{
		"t1": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         2,
			y:         3,
			// true because v[2] = 0xff, and v[3] = 0xff, so an increment of 2 should occur
			wantedInc: true,
		},
		"t2": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0xff, 0x22, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         3,
			y:         4,
			// false because v[2] = 0xff, and v[3] = 0x22, so an increment of 2 should not occur
			wantedInc: false,
		},
		"t3": {
			pc: 0x0333,
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         18,
			y:         4,
			wantedErr: errors.New("register index out of bounds"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.pc = test.pc
			em.registers = test.registers
			err := em.seqVxVy(test.x, test.y)

			assert.True(t, (err != nil) == (test.wantedErr != nil))
			if test.wantedErr != nil {
				assert.Equal(t, test.wantedErr, err)
				return
			}
			if test.wantedInc {
				assert.Equal(t, test.pc+2, em.pc)
			} else {
				assert.Equal(t, test.pc, em.pc)
			}
		})
	}
}

func Test_ldVxKK(t *testing.T) {
	tests := map[string]struct {
		registers []uint8
		x         uint16
		kk        uint16
		wantedErr error
	}{
		"t1": {
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         2,
			kk:        0x03,
		},
		"t2": {
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         18,
			kk:        0x03,
			wantedErr: errors.New("register index out of bounds"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.registers = test.registers
			err := em.ldVxKK(test.x, test.kk)

			assert.True(t, (err != nil) == (test.wantedErr != nil))
			if test.wantedErr != nil {
				assert.Equal(t, test.wantedErr, err)
				return
			}
			assert.Equal(t, em.registers[test.x], uint8(test.kk))
		})
	}
}

func Test_addVxKK(t *testing.T) {
	tests := map[string]struct {
		registers []uint8
		x         uint16
		kk        uint16
		wantedErr error
	}{
		"t1": {
			// v
			registers: []uint8{0x21, 0x33, 0xee, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         2,
			kk:        0x03,
		},
		"t2": {
			// v
			registers: []uint8{0x21, 0x33, 0xff, 0x12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			x:         18,
			kk:        0x03,
			wantedErr: errors.New("register index out of bounds"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			em := Create()
			em.registers = test.registers
			var before uint8 = 0
			if test.wantedErr == nil {
				before = em.registers[test.x]
			}
			err := em.addVxKK(test.x, test.kk)

			assert.True(t, (err != nil) == (test.wantedErr != nil))
			if test.wantedErr != nil {
				assert.Equal(t, test.wantedErr, err)
				return
			}
			assert.Equal(t, before+uint8(test.kk), em.registers[test.x])
		})
	}
}
