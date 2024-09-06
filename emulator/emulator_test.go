package emulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetch(t *testing.T) {
	tests := map[string]struct {
		pc     uint16
		p1, p2 uint8
		inst   uint16
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

			inst := em.fetch()
			assert.Equal(t, test.inst, inst)
			assert.Equal(t, test.pc+2, em.pc)
		})
	}
}
