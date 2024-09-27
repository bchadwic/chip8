package emulator

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bchadwic/chip8/internal/display"
	"github.com/bchadwic/chip8/internal/display/emit"
	"github.com/bchadwic/chip8/internal/drivers"
	"github.com/bchadwic/chip8/internal/keypad"
	"github.com/bchadwic/chip8/internal/speaker"
)

const (
	REGISTERS  = 16
	MEM_SIZE   = 4096
	STACK_SIZE = 16

	ROWS, COLS = 32, 64

	FONT_ADDR = 0x050
	ROM_ADDR  = 0x200

	// niblet masks
	N1_MASK = 0xF000
	N2_MASK = 0x0F00
	N3_MASK = 0x00F0
	N4_MASK = 0x000F

	CLS_OR_RET = 0x0000

	// instructions
	CLS           = 0x00E0 // clear screen
	RET           = 0x00EE // return from subroutine
	JMP           = 0x1000 // jump pc to address
	CALL          = 0x2000 // call subroutine
	SEQ_VX_NN     = 0x3000 // skip if vx eq nn
	SNE_VX_NN     = 0x4000 // skip if vx ne nn
	SEQ_VX_VY     = 0x5000 // skip if vx eq vy
	LD_VX_KK      = 0x6000 // load vx with kk
	ADD_VX_KK     = 0x7000 // load vx with kk
	MOD_VX_VY_OPS = 0x8000 // series of arithmetic and bit operations for vx / vy
	SNE_VX_VY     = 0x9000 // skip if vx ne vy
	LD_I          = 0xA000 // load register i with remaining bits
	JMP_V0        = 0xB000 // jump to nnn + v0
	RND_VX_KK     = 0xC000 // generate random, then bitwise and kk, store to vx
	DRW_VX_VY_N   = 0xD000 // generate random, then bitwise and kk, store to vx
	VX_KEY_OPS    = 0xE000 // series of skip instructions for key presses
	TIMING_OPS    = 0xF000 // series of timing instructions

	// sub instructions under MOD_VX_VY_OPS
	LD_VX_VY   = 0x0000 // store vx in vy
	OR_VX_VY   = 0x0001 // bitwise vx or vy
	AND_VX_VY  = 0x0002 // bitwise vx and vy
	XOR_VX_VY  = 0x0003 // bitwise vx xor vy
	ADD_VX_VY  = 0x0004 // vx + vy
	SUB_VX_VY  = 0x0005 // vx - vy
	SHR_VX_VY  = 0x0006 // vx >> 1 and store lsb in vf
	SUBN_VX_VY = 0x0007 // vx - vy and set underflow in vf
	SHL_VX_VY  = 0x000E // vx << 1 and store msb in vf

	// sub instructions under VX_KEY_OPS
	SEQ_VX_KEY_PR = 0x009E // skip if vx contains key pressed
	SNE_VX_KEY_PR = 0x00A1 // skip if vx does not contain key pressed

	// sub instructions under TIMING_OPS
	LD_VX_DT = 0x0007 // set vx to delay timer
	LD_VX_K  = 0x000A // wait for key press, store in vx
	LD_DT_VX = 0x0015 // set delay timer to vx
	LD_ST_VX = 0x0018 // set sound timer to vx
	ADD_I_VX = 0x001E // i + vx and set to i
	LD_F_VX  = 0x0029 // set i to sprite stored in vx
	LD_B_VX  = 0x0033 // i, i+1, and i+2 represent BCD of vx
	LD_I_VX  = 0x0055 // load memory i-n with the values stored in v0-vx
	LD_VX_I  = 0x0065 // load v0-vx with values stored in memory i-n
)

var fonts []uint8 = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type EmulatorSettings struct {
	FrameRate int
	Rom       []uint8
	Fill      bool
	Color     string
	Keyboard  string
}

type emulator struct {
	registers []uint8
	mem       []uint8

	// stack pointer
	sp    uint8
	stack []uint16

	// index register
	i uint16

	// program counter
	pc uint16

	// delay and sound timers
	dt, st uint8

	settings *EmulatorSettings

	// devices
	speaker speaker.Speaker
	keypad  keypad.Keypad
	display display.Display
}

func Create(settings *EmulatorSettings) *emulator {
	mem := make([]uint8, MEM_SIZE)
	// load fonts into memory
	for i := 0; i < len(fonts); i++ {
		mem[i+FONT_ADDR] = fonts[i]
	}

	speaker := speaker.Create()
	keypad := keypad.Create()
	display := display.Create(ROWS, COLS)

	go drivers.Create(
		speaker,
		keypad,
		display,
	).KeypadSettings(
		settings.Keyboard,
	).DisplaySettings(
		settings.FrameRate,
		settings.Fill,
		settings.Color,
	).Start()

	return &emulator{
		registers: make([]uint8, REGISTERS),
		mem:       mem,
		stack:     make([]uint16, STACK_SIZE),
		settings:  settings,
		speaker:   speaker,
		keypad:    keypad,
		display:   display,
	}
}

func (em *emulator) Load(rom []uint8) {
	for i := 0; i < len(rom); i++ {
		em.mem[i+ROM_ADDR] = rom[i]
	}
	em.pc = ROM_ADDR
}

func (em *emulator) Start() {
	clock := time.NewTicker(time.Duration(em.settings.FrameRate) * time.Millisecond)
	defer clock.Stop()

	i := 0
	for range clock.C {
		em.speaker.Set(em.st > 0)
		if em.speaker.IsActive() {
			em.st--
		}
		if em.dt > 0 {
			em.dt--
		}
		inst, err := em.fetch()
		if err != nil {
			log.Fatal(err.Error())
		}
		em.execute(inst)
		if i%10 == 0 {
			em.keypad.Clear()
		}
		i++
	}
}

// fetch retrieves two bytes located at pc
// if two bytes are not available within
// the available memory, error is returned
func (em *emulator) fetch() (uint16, error) {
	if em.pc > MEM_SIZE-2 {
		return 0, fmt.Errorf("pc out of memory bounds: %d", em.pc)
	}
	p1 := em.mem[em.pc]
	p2 := em.mem[em.pc+1]
	return (uint16(p1) << 8) | uint16(p2), nil
}

func (em *emulator) execute(inst uint16) {
	n1 := inst & N1_MASK
	n2 := inst & N2_MASK
	n3 := inst & N3_MASK
	n4 := inst & N4_MASK

	addr := n2 | n3 | n4
	x := n2 >> 8
	y := n3 >> 4
	nn := n3 | n4

	inc := true
	switch n1 {
	case CLS_OR_RET:
		switch inst {
		case CLS:
			em.cls()
		case RET:
			em.ret()
		default:
			log.Fatalf("opcode not found: %x", inst)
		}
	case JMP:
		em.jmp(addr)
		inc = false
	case CALL:
		em.call(addr)
		inc = false
	case SEQ_VX_NN:
		em.seqVxNN(x, nn)
	case SNE_VX_NN:
		em.sneVxNN(x, nn)
	case SEQ_VX_VY:
		if n4 == 0 {
			em.seqVxVy(x, y)
		} else {
			log.Fatalf("opcode not found: %x", inst)
		}
	case LD_VX_KK:
		em.ldVxKK(x, nn)
	case ADD_VX_KK:
		em.addVxKK(x, nn)
	// TODO test these functions
	case MOD_VX_VY_OPS:
		switch n4 {
		case LD_VX_VY:
			em.ldVxVy(x, y)
		case OR_VX_VY:
			em.orVxVy(x, y)
		case AND_VX_VY:
			em.andVxVy(x, y)
		case XOR_VX_VY:
			em.xorVxVy(x, y)
		case ADD_VX_VY:
			em.addVxVy(x, y)
		case SUB_VX_VY:
			em.subVxVy(x, y)
		case SHR_VX_VY:
			em.shrVxVy(x, y)
		case SUBN_VX_VY:
			em.subnVxVy(x, y)
		case SHL_VX_VY:
			em.shlVxVy(x, y)
		default:
			log.Fatalf("opcode not found: %x", inst)
		}
	case SNE_VX_VY:
		if n4 == 0 {
			em.sneVxVy(x, y)
		} else {
			log.Fatalf("opcode not found: %x", inst)
		}
	case LD_I:
		em.ldI(addr)
	case JMP_V0:
		em.jmpV0(addr)
		inc = false
	case RND_VX_KK:
		em.rndVxKK(x, nn)
	case DRW_VX_VY_N:
		em.drawVxVyN(x, y, n4)
	case VX_KEY_OPS:
		if nn == SEQ_VX_KEY_PR {
			em.seqVxKey(x)
		} else if n3|n4 == SNE_VX_KEY_PR {
			em.sneVxKey(x)
		} else {
			log.Fatalf("opcode not found: %x", inst)
		}
	case TIMING_OPS:
		switch nn {
		case LD_VX_DT:
			em.ldVxDt(x)
		case LD_VX_K:
			em.ldVxK(x)
		case LD_DT_VX:
			em.ldDtVx(x)
		case LD_ST_VX:
			em.ldStVx(x)
		case ADD_I_VX:
			em.addIVx(x)
		case LD_F_VX:
			em.ldFVx(x)
		case LD_B_VX:
			em.ldBVx(x)
		case LD_I_VX:
			em.ldIVx(x)
		case LD_VX_I:
			em.ldVxI(x)
		default:
			log.Fatalf("opcode not found: %x", inst)
		}
	}
	if inc {
		em.pc += 2
	}
}

// clear screen
func (em *emulator) cls() {
	em.display.Clear()
}

// return from subroutine
func (em *emulator) ret() {
	em.sp--
	em.pc = em.stack[em.sp]
}

// jump program counter to instructed address
func (em *emulator) jmp(addr uint16) {
	em.pc = addr
}

// call subroutine
func (em *emulator) call(addr uint16) {
	// move stack pointer to next position, save current position of program counter
	em.stack[em.sp] = em.pc
	em.sp++
	em.pc = addr
}

// 0x3XNN
// skip to next instruction set if register X is equal to NN
func (em *emulator) seqVxNN(x uint16, nn uint16) {
	if em.registers[x] == uint8(nn) {
		em.pc += 2
	}
}

// 0x4XNN
// skip to next instruction set if register X is NOT equal to NN
func (em *emulator) sneVxNN(x uint16, nn uint16) {
	if em.registers[x] != uint8(nn) {
		em.pc += 2
	}
}

// 0x5XY0
// skip to next instruction set if register X is equal to register Y
func (em *emulator) seqVxVy(x uint16, y uint16) {
	if em.registers[x] == em.registers[y] {
		em.pc += 2
	}
}

// 0x6XKK
// load register X with the value of KK
func (em *emulator) ldVxKK(x uint16, kk uint16) {
	em.registers[x] = uint8(kk)
}

// 0x7XKK
// add the value of KK to register X
func (em *emulator) addVxKK(x uint16, kk uint16) {
	em.registers[x] += uint8(kk)
}

// 0x8xy0
// store register X value in register Y
func (em *emulator) ldVxVy(x uint16, y uint16) {
	em.registers[x] = em.registers[y]
}

// 0x8xy1
// bitwise register X or Y, then store to register X
func (em *emulator) orVxVy(x uint16, y uint16) {
	em.registers[x] |= em.registers[y]
}

// 0x8xy2
// bitwise register X and Y, then store to register X
func (em *emulator) andVxVy(x uint16, y uint16) {
	em.registers[x] &= em.registers[y]
}

// 0x8xy3
// bitwise register X xor Y, then store to register X
func (em *emulator) xorVxVy(x uint16, y uint16) {
	em.registers[x] ^= em.registers[y]
}

// 0x8xy4
// add register X and Y, then store to register X
// if overflow occurs, set VF register to 1
func (em *emulator) addVxVy(x uint16, y uint16) {
	sum := em.registers[x] + em.registers[y]
	if sum < em.registers[x] {
		em.registers[0xF] = 1 // set overflow
	}
	em.registers[x] = sum
}

// 0x8xy5
// subtract register Y from X, then store to register X
// if underflow occurs, set VF register to 0, otherwise 1
func (em *emulator) subVxVy(x uint16, y uint16) {
	diff := em.registers[x] - em.registers[y]
	if em.registers[y] > em.registers[x] {
		em.registers[0xF] = 0
	} else {
		em.registers[0xF] = 1 // set underflow
	}
	em.registers[x] = diff
}

// 0x8xy6
// store the LSB of the value stored in register X to VF
// then right shift the value of register X by 1, then store to register X
func (em *emulator) shrVxVy(x uint16, _ uint16) {
	em.registers[0xF] = em.registers[x] & 0x01
	em.registers[x] >>= 1
}

// 0x8xy7
// subtract register X from Y, then store to register X
// if underflow occurs, set VF register to 0, otherwise 1
func (em *emulator) subnVxVy(x uint16, y uint16) {
	if em.registers[x] > em.registers[y] {
		em.registers[0xF] = 0
	} else {
		em.registers[0xF] = 1
	}
	em.registers[x] = em.registers[y] - em.registers[x]
}

// 0x8xyE
// store the MSB of the value stored in register X to VF
// then left shift the value of register X by 1, then store to register X
func (em *emulator) shlVxVy(x uint16, _ uint16) {
	em.registers[0xF] = em.registers[x] >> 7
	em.registers[x] <<= 1
}

// 0x9xy0
// skip to next instruction set if register X is NOT equal to register Y
func (em *emulator) sneVxVy(x uint16, y uint16) {
	if em.registers[x] != em.registers[y] {
		em.pc += 2
	}
}

// 0xAnnn
// set the value of register I to addr (nnn)
func (em *emulator) ldI(addr uint16) {
	em.i = addr
}

// 0xBnnn
// set the program counter to addr (nnn) + register v0 value
func (em *emulator) jmpV0(addr uint16) {
	em.pc = addr + uint16(em.registers[0])
}

// 0xCxkk
// set register X to the value of a random number bitwise and KK
func (em *emulator) rndVxKK(x uint16, kk uint16) {
	r := uint8(rand.Intn(0x100))
	em.registers[x] = r & uint8(kk)
}

// 0xDxyn
// draw a sprite at register X and Y location, of N height
func (em *emulator) drawVxVyN(x uint16, y uint16, n uint16) {
	startc := em.registers[x] % COLS // clamp cx to display width
	startr := em.registers[y] % ROWS // clamp cy to display height
	em.registers[0xF] = 0            // clear collision flag

	for rowi := uint8(0); rowi < uint8(n); rowi++ {
		index := em.i + uint16(rowi)
		row := em.mem[index]
		// loop through each bit in a byte (8)
		for coli := uint8(0); coli < 8; coli++ {
			// read sprite bits left to right
			spb := row & (0b1000_0000 >> coli)
			if spb == 0 {
				continue
			}

			pixelr := (startr + rowi) % COLS
			pixelc := startc + coli
			if startc+coli < COLS && startr+rowi < ROWS {
				pixel := em.display.Get(pixelr, pixelc)
				if pixel {
					em.display.Set(emit.OFF, pixelr, pixelc)
					em.registers[0xF] = 1 // Collision detected
				} else {
					em.display.Set(emit.ON, pixelr, pixelc)
				}
			}
		}
	}
}

// 0xEX9E
// skip the next instruction if the key at register X was pressed
func (em *emulator) seqVxKey(x uint16) {
	key := em.registers[x]
	// check if key is pressed
	if em.keypad.Get(key) {
		em.pc += 2
	}
}

// 0xEXA1
// skip the next instruction if the key at register X was not pressed
func (em *emulator) sneVxKey(x uint16) {
	key := em.registers[x]
	// check if key is not pressed
	if !em.keypad.Get(key) {
		em.pc += 2
	}
}

// 0xFX07
// set the value of register X to delay timer
func (em *emulator) ldVxDt(x uint16) {
	em.registers[x] = em.dt
}

// 0xFX0A
// await a keypress, and assign keycode to register X
func (em *emulator) ldVxK(x uint16) {
	kaddr := em.keypad.Next()
	em.registers[x] = kaddr
}

// 0xFX15
// set the delay timer to the value of register X
func (em *emulator) ldDtVx(x uint16) {
	em.dt = em.registers[x]
}

// 0xFX18
// set the sound timer to the value of register X
func (em *emulator) ldStVx(x uint16) {
	em.st = em.registers[x]
}

// 0xFX1E
// add i and value of register X, then store to i
func (em *emulator) addIVx(x uint16) {
	if em.i+uint16(em.registers[x]) > 0xFFF {
		em.registers[0xF] = 1
	} else {
		em.registers[0xF] = 0
	}
	em.i += uint16(em.registers[x])
}

// 0xFX29
// add i and value of register X, then store to i
func (em *emulator) ldFVx(x uint16) {
	em.i = uint16(em.registers[x]) * 5
}

// 0xFX33
// store BCD representation of the value stored in register X
// in memory locations I, I+1, and I+2.
func (em *emulator) ldBVx(x uint16) {
	bcd := em.registers[x]
	least := bcd % 10
	mid := ((bcd % 100) - least) / 10
	most := ((bcd) - ((bcd % 100) - least) - least) / 100

	// write to memory
	em.mem[em.i] = most
	em.mem[em.i+1] = mid
	em.mem[em.i+2] = least
}

// 0xFX55
// store the values in registers 0-X to memory starting at i
func (em *emulator) ldIVx(x uint16) {
	for i := uint16(0); i <= x; i++ {
		em.mem[em.i+i] = em.registers[i]
	}
	em.i = x + 1
}

// 0xFX65
// store the values in memory starting at i into registers 0-X
func (em *emulator) ldVxI(x uint16) {
	for i := uint16(0); i <= x; i++ {
		em.registers[i] = em.mem[em.i+i]
	}
	em.i = x + 1
}
