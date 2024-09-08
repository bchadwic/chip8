package keypad

type Keypad interface {
	// checks if the key addr passed in is in the down position
	IsPressed(kaddr uint8) bool
	// blocking call, awaits the next key press
	GetNextKey() uint8
}
