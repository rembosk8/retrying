package retrying

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetResetDefault(t *testing.T) {
	dSettings := newSettings()

	SetDefault(WithDuration(10 * time.Hour))

	newDefSettings := newSettings()
	assert.Equal(t, 10*time.Hour, newDefSettings.Duration)
	assert.NotEqual(t, dSettings, newDefSettings)

	ResetDefault()

	resetSettings := newSettings()
	assert.Equal(t, dSettings, resetSettings)
}
