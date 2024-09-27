package keypad

import (
	"sync"
)

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
	kp.pressed[kaddr] = true
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
