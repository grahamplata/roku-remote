package device

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

const helpText = `Roku Remote Control
q: Quit
p: Power
+/-: Volume up/down
m: Mute
Arrow keys: Navigate
Enter: Select
b: Back
h: Home
r: Rewind
f: Fast Forward
Space: Play
i: Instant Replay
Tab: Info
Backspace: Backspace
/: Search
Ctrl+F: Find Remote
Page Up: Channel Up
Page Down: Channel Down
t: Tuner
1-4: HDMI 1-4`

var keyCommands = map[string]string{
	"p":         "poweroff",
	"+":         "volumeup",
	"-":         "volumedown",
	"m":         "mute",
	"up":        "up",
	"down":      "down",
	"left":      "left",
	"right":     "right",
	"enter":     "select",
	"b":         "back",
	"h":         "home",
	"r":         "rev",
	"f":         "fwd",
	" ":         "play",
	"i":         "replay",
	"tab":       "info",
	"backspace": "backspace",
	"/":         "search",
	"ctrl+f":    "find",
	"pgup":      "channelup",
	"pgdown":    "channeldown",
	"t":         "tuner",
	"1":         "HDMI1",
	"2":         "HDMI2",
	"3":         "HDMI3",
	"4":         "HDMI4",
}

func ControlCmd(ch *cmdutil.Helper) *cobra.Command {
	var controlCmd = &cobra.Command{
		Use:   "control",
		Short: "Control a Roku device via keyboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				return fmt.Errorf("invalid Roku host: %w", err)
			}
			device := roku.NewDevice(ip)
			p := tea.NewProgram(&controlModel{device: device, ctx: ctx})
			if _, err := p.Run(); err != nil {
				return err
			}
			return nil
		},
	}
	return controlCmd
}

// controlModel represents the state of the device control interface
type controlModel struct {
	// ctx is the context for device actions.
	ctx context.Context
	// device represents the target device to control.
	device *roku.Device
	// statusMessage displays feedback to the user.
	statusMessage string
}

// Init initializes the model
func (m controlModel) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model accordingly.
func (m *controlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for context cancellation to handle graceful shutdown
	select {
	case <-m.ctx.Done():
		return m, tea.Quit
	default:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case errorMsg:
		m.statusMessage = fmt.Sprintf("Error sending command '%s': %v", msg.cmd, msg.err)
		return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg { return clearStatusMsg{} })
	case successMsg:
		m.statusMessage = fmt.Sprintf("Command '%s' sent", msg.cmd)
		return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg { return clearStatusMsg{} })
	case clearStatusMsg:
		m.statusMessage = ""
	}
	return m, nil
}

// clearStatusMsg is a message to clear the status message.
type clearStatusMsg struct{}

// errorMsg represents an error from sending a command.
type errorMsg struct {
	err error
	cmd string
}

// successMsg represents a successful command send.
type successMsg struct {
	cmd string
}

func (m controlModel) View() string {
	return fmt.Sprintf("%s\n\n%s", m.statusMessage, helpText)
}

// handleKeyPress processes a key press and sends the corresponding command to the Roku device.
func (m *controlModel) handleKeyPress(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit
		if key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))) {
			return m, tea.Quit
		}

		// Map key presses to Roku commands
		if cmd, exists := keyCommands[msg.String()]; exists {
			m.statusMessage = fmt.Sprintf("Sending '%s'...", cmd)
			return m, m.sendCommand(cmd)
		}

		// Log unknown keys for debugging
		log.Printf("Unknown key pressed: %s", msg.String())
	}
	return m, nil
}

// sendCommand sends a command to the device asynchronously and updates the status.
func (m *controlModel) sendCommand(cmd string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := m.device.Action(m.ctx, cmd)
		if err != nil {
			return errorMsg{err, cmd}
		}
		return successMsg{cmd}
	})
}
