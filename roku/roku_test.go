package roku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvailableActions(t *testing.T) {
	actions := AvailableActions()

	// Verify the map is not empty
	assert.NotEmpty(t, actions)

	// Verify some known actions exist
	expectedActions := map[string]string{
		"home":       "/Home",
		"select":     "/Select",
		"up":         "/Up",
		"down":       "/Down",
		"left":       "/Left",
		"right":      "/Right",
		"volumeup":   "/VolumeUp",
		"volumedown": "/VolumeDown",
		"mute":       "/VolumeMute",
		"poweroff":   "/PowerOff",
		"HDMI1":      "/InputHDMI1",
	}

	for action, expectedPath := range expectedActions {
		actualPath, exists := actions[action]
		assert.True(t, exists, "Expected action %q to exist", action)
		assert.Equal(t, expectedPath, actualPath, "Action %q should map to %q", action, expectedPath)
	}

	// Verify total count (as of current implementation, there are 27 actions)
	assert.Equal(t, 27, len(actions), "Expected 27 total actions")
}
