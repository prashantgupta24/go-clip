package systray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Init(t *testing.T) {
	assert.NotNil(t, clipboardInstance)
	assert.Equal(t, clipboardInstance.activeSlots, 20)
}
