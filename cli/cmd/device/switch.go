package device

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SwitchCmd(ch *cmdutil.Helper) *cobra.Command {
	var switchCmd = &cobra.Command{
		Use:   "switch",
		Short: "Switch the default Roku device.",
		Long:  `Select a different Roku device from the stored list to set as the default.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			devices := viper.GetStringSlice("roku.devices")
			if len(devices) == 0 {
				fmt.Println("No devices stored. Run 'roku find' to discover and store devices.")
				return nil
			}
			return runSwitch(ctx, devices)
		},
	}

	return switchCmd
}

func runSwitch(ctx context.Context, devices []string) error {
	p := tea.NewProgram(initialSwitchModel(ctx, devices))
	m, err := p.Run()
	if err != nil {
		return err
	}

	finalModel, ok := m.(switchModel)
	if !ok {
		return fmt.Errorf("unexpected model type returned from tea.Program")
	}
	if finalModel.selected >= 0 {
		selectedIP := devices[finalModel.selected]
		viper.Set("roku.host", selectedIP)
		if err := AddToConfigFile(); err != nil {
			return fmt.Errorf("failed to save device configuration: %w", err)
		}
	}
	return nil
}

type switchModel struct {
	devices  []string
	cursor   int
	selected int
	ctx      context.Context
}

func initialSwitchModel(ctx context.Context, devices []string) switchModel {
	return switchModel{
		devices:  devices,
		cursor:   0,
		selected: -1,
		ctx:      ctx,
	}
}

func (m switchModel) Init() tea.Cmd {
	return nil
}

func (m switchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for context cancellation
	select {
	case <-m.ctx.Done():
		return m, tea.Quit
	default:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.devices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m switchModel) View() string {
	s := "Select a Roku device to set as default:\n\n"

	for i, device := range m.devices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, device)
	}

	s += "\nPress q to quit, enter to select.\n"
	return s
}
