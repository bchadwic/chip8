package keypad

/*

package keypad

import (
	"sync"
)

// dvorak to keypad
var toKeypad map[byte]uint8 = map[byte]uint8{
	'1': 0x1, '2': 0x2, '3': 0x3, '4': 0xC,
	'\'': 0x4, ',': 0x5, '.': 0x6, 'p': 0xD,
	'a': 0x7, 'o': 0x8, 'e': 0x9, 'u': 0xE,
	';': 0xA, 'q': 0x0, 'j': 0xB, 'k': 0xF,
}

type Keypad interface {
	Clear()
	Get(kaddr uint8) bool
	Set(kaddr uint8)
	Next() uint8
}

type keypad struct {
	pressed map[byte]bool
	mu      sync.Mutex
}

func Create() Keypad {
	return &keypad{
		pressed: make(map[byte]bool),
	}
}

func (kp *keypad) Clear() {
	kp.mu.Lock()
	defer kp.mu.Unlock()
	for kaddr := range kp.pressed {
		kp.pressed[kaddr] = false
	}
}

func (kp *keypad) Get(kaddr uint8) bool {
	kp.mu.Lock()
	defer kp.mu.Unlock()
	return kp.pressed[kaddr]
}

func (kp *keypad) Set(kaddr uint8) {
	kp.mu.Lock()
	defer kp.mu.Unlock()
	kp.pressed[toKeypad[kaddr]] = true
}

func (kp *keypad) Next() uint8 {
	for {
		kp.mu.Lock()
		for kaddr, pressed := range kp.pressed {
			if pressed {
				kp.mu.Unlock()
				return kaddr
			}
		}
		kp.mu.Unlock()
	}
}

*/
