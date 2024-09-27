package keypad

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	keypad := Create()
	assert.NotNil(t, keypad)
}

func Test_Clear(t *testing.T) {
	keypad := &keypad{
		pressed: map[byte]bool{
			'a': true,
		},
	}
	keypad.Clear()
	assert.False(t, keypad.pressed['a'])
}

func Test_Get(t *testing.T) {
	keypad := &keypad{
		pressed: map[byte]bool{
			'a': true,
		},
	}
	assert.True(t, keypad.Get('a'))
}

func Test_Set(t *testing.T) {
	keypad := &keypad{
		pressed: map[byte]bool{},
	}
	keypad.Set('a')
	assert.True(t, keypad.pressed['a'])
}

func Test_Next(t *testing.T) {
	keypad := &keypad{
		pressed: map[byte]bool{},
	}
	go func() {
		assert.Equal(t, uint8('a'), keypad.Next())
	}()
	time.Sleep(2 * time.Second)
	keypad.Set('a')
}
