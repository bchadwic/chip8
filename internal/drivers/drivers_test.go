package drivers

import (
	"testing"

	"github.com/gonutz/prototype/draw"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	dc := Create(nil, nil, nil)
	assert.NotNil(t, dc)
}

func Test_DisplaySettings(t *testing.T) {
	dc := Create(nil, nil, nil)
	dc.DisplaySettings(1, true, "GrAy")
	assert.True(t, dc.displayInitialized)
	assert.True(t, dc.fill)
	assert.Equal(t, 1, dc.frameRate)
	assert.Equal(t, draw.Gray, dc.color)
}

func Test_KeypadSettings(t *testing.T) {
	dc := Create(nil, nil, nil)
	dc.KeypadSettings("Qwerty")
	assert.True(t, dc.keypadInitialized)
	assert.Equal(t, qwerty, dc.keyboard)
}
