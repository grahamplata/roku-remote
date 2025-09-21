package api

// actions - the available actions for the roku device
// - Docs: https://developer.roku.com/docs/developer-program/debugging/external-control-api.md#keypress-key-values
var ExternalControlActions = map[string]string{
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
