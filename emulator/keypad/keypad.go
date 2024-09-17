package keypad

import (
	"fmt"
	"sync"
)

type Keypad interface {
	Clear()
	Get(kaddr uint8) bool
	Set(kaddr uint8)
	Next() uint8
}

type keypad struct {
	keys    map[byte]bool
	pressed map[byte]bool
	ikeys   []byte
	mu      sync.Mutex
}

func Create(keys string) Keypad {
	mkeys := make(map[byte]bool, len(keys))
	mpressed := make(map[byte]bool, len(keys))
	ikeys := make([]byte, len(keys))
	for i, k := range keys {
		bk := byte(k)
		mkeys[bk] = true
		ikeys[i] = bk
	}
	return &keypad{
		keys:    mkeys,
		pressed: mpressed,
		ikeys:   ikeys,
	}
}

func (kp *keypad) Clear() {
	kp.mu.Lock()
	defer kp.mu.Unlock()
	for kaddr := range kp.keys {
		kp.pressed[kaddr] = false
	}
}

func (kp *keypad) Get(kaddr uint8) bool {
	// fmt.Printf("getting key: %c\n", kp.ikeys[kaddr])
	kp.mu.Lock()
	defer kp.mu.Unlock()
	return kp.pressed[kp.ikeys[kaddr]]
}

func (kp *keypad) Set(kaddr uint8) {
	kp.mu.Lock()
	defer kp.mu.Unlock()
	if !kp.keys[kaddr] {
		fmt.Printf("invalid key: %c\n", kaddr)
		return
	}
	// fmt.Printf("valid key: %c\n", kaddr)
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
