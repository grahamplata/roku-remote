package roku

import (
	"sort"
	"strings"
)

// endpoints - the available endpoints for the roku device
// - Docs: https://developer.roku.com/docs/developer-program/debugging/external-control-api.md#general-ecp-commands
var endpoints = map[string]string{
	"apps":     "/query/apps",
	"active":   "/query/active-app",
	"player":   "/query/media-player",
	"device":   "/query/device-info",
	"icon":     "/query/icon",
	"input":    "/input",
	"search":   "/search",
	"keydown":  "/keydown",
	"keypress": "/keypress",
	"launch":   "/launch",
	"install":  "/install",
}

// actions - the available actions for the roku device
// - Docs: https://developer.roku.com/docs/developer-program/debugging/external-control-api.md#keypress-key-values
var actions = map[string]string{
	"home":        "/Home",
	"rev":         "/Rev",
	"fwd":         "/Fwd",
	"play":        "/Play",
	"select":      "/Select",
	"left":        "/Left",
	"right":       "/Right",
	"down":        "/Down",
	"up":          "/Up",
	"back":        "/Back",
	"replay":      "/InstantReplay",
	"info":        "/Info",
	"backspace":   "/Backspace",
	"search":      "/Search",
	"enter":       "/Enter",
	"find":        "/FindRemote",
	"volumedown":  "/VolumeDown",
	"mute":        "/VolumeMute",
	"volumeup":    "/VolumeUp",
	"poweroff":    "/PowerOff",
	"channelup":   "/ChannelUp",
	"channeldown": "/ChannelDown",
	"tuner":       "/InputTuner",
	"HDMI1":       "/InputHDMI1",
	"HDMI2":       "/InputHDMI2",
	"HDMI3":       "/InputHDMI3",
	"HDMI4":       "/InputHDMI4",
}

// AvialableActions returns the available action commands
func AvialableActions() string {
	keys := make([]string, 0, len(actions))
	for k := range actions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
