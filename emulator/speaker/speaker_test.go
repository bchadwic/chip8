package speaker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	speaker := Create()
	assert.NotNil(t, speaker)
}

func Test_IsActive(t *testing.T) {
	speaker := &speaker{}
	assert.False(t, speaker.IsActive())
	speaker.active = true
	assert.True(t, speaker.IsActive())
}

func Test_Set(t *testing.T) {
	speaker := Create()
	speaker.Set(true)
	assert.True(t, speaker.IsActive())
}
