package device

import (
	"context"
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

func SendCmd(ch *cmdutil.Helper) *cobra.Command {
	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send an action to your Roku Device.",
		Long:  "Interactively select and send actions to your Roku device.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				return fmt.Errorf("invalid Roku host: %w", err)
			}
			return runSend(ctx, ip)
		},
	}
	return sendCmd
}

func runSend(ctx context.Context, ip string) error {
	actions := roku.AvailableActions()
	var actionNames []string
	for name := range actions {
		actionNames = append(actionNames, name)
	}
	sort.Strings(actionNames)

	p := tea.NewProgram(initialSendModel(actionNames, ctx, ip))
	m, err := p.Run()
	if err != nil {
		return err
	}

	finalModel := m.(sendModel)
	if finalModel.selected >= 0 {
		selectedAction := actionNames[finalModel.selected]
		device := roku.NewDevice(ip)
		err := device.Action(ctx, selectedAction)
		if err != nil {
			return fmt.Errorf("error sending action: %w", err)
		}
		fmt.Println("Action sent successfully.")
	}
	return nil
}

type sendModel struct {
	actions  []string
	cursor   int
	selected int
	ctx      context.Context
	ip       string
}

func initialSendModel(actions []string, ctx context.Context, ip string) sendModel {
	return sendModel{
		actions:  actions,
		cursor:   0,
		selected: -1,
		ctx:      ctx,
		ip:       ip,
	}
}

func (m sendModel) Init() tea.Cmd {
	return nil
}

func (m sendModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.actions)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m sendModel) View() string {
	s := "Select an action to send to your Roku device:\n\n"

	for i, action := range m.actions {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, action)
	}

	s += "\nPress q to quit, enter to send.\n"
	return s
}
