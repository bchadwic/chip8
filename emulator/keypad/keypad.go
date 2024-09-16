package keypad

import "fmt"

type Keypad interface {
	// checks if the key addr passed in is in the down position
	IsPressed(kaddr uint8) bool
	IsPressedFunc(f func(kaddr uint8) bool)
	// blocking call, awaits the next key press
	GetNextKey() uint8
	GetNextKeyFunc(f func() uint8)
	// continuous stream of keys being pressed
	KeyStream() chan uint8
	// representation of keys pressed
	KeyMap() map[uint8]bool
}

var (
	ValidKeys = map[uint8]bool{
		'1': true, '2': true, '3': true, '4': true,
		'\'': true, ',': true, '.': true, 'p': true,
		'a': true, 'o': true, 'e': true, 'u': true,
		';': true, 'x': true, 'c': true, 'v': true,
	}
)

type keypad struct {
	stream     chan uint8
	keys       map[uint8]bool
	isPressed  func(kaddr uint8) bool
	getNextKey func() uint8
}

func Create() *keypad {
	keys := map[uint8]bool{
		'1': false, '2': false, '3': false, '4': false,
		'\'': false, ',': false, '.': false, 'p': false,
		'a': false, 'o': false, 'e': false, 'u': false,
		';': false, 'x': false, 'c': false, 'v': false,
	}
	stream := make(chan uint8)
	return &keypad{
		keys:   keys,
		stream: stream,
	}
}

func (kp *keypad) IsPressed(kaddr uint8) bool {
	if kp.isPressed == nil {
		fmt.Println("is pressed is nil")
		return false
	}
	return kp.isPressed(kaddr)
}

func (kp *keypad) IsPressedFunc(f func(kaddr uint8) bool) {
	kp.isPressed = f
}

func (kp *keypad) GetNextKey() uint8 {
	if kp.getNextKey == nil {
		fmt.Println("get next key is nil")
		return 0
	}
	return kp.getNextKey()
}

func (kp *keypad) GetNextKeyFunc(f func() uint8) {
	kp.getNextKey = f
}

func (kp *keypad) KeyMap() map[uint8]bool {
	return kp.keys
}

func (kp *keypad) KeyStream() chan uint8 {
	return kp.stream
}
