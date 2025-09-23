package device

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DefaultScanTime = 5

func FindCmd(ch *cmdutil.Helper) *cobra.Command {
	var findCmd = &cobra.Command{
		Use:   "find",
		Short: "Find Roku Remotes on your local network.",
		Long: `Scans local network using SSDP to locate Roku 
devices that you can interface with.
    
This command uses Simple Service Discovery Protocol (or SSDP) which
provides a mechanism where by network clients, with little or no
static configuration, can discover network services.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			wait, err := cmd.Flags().GetInt("wait")
			if err != nil {
				return fmt.Errorf("unable to complete (find) command: %w", err)
			}
			fmt.Printf("Scanning for Roku devices for %d seconds...\n", wait)
			devices, err := roku.Find(wait)
			if err != nil {
				return fmt.Errorf("unable to complete (find) command: %w", err)
			}
			return runFind(devices)
		},
	}
	findCmd.Flags().IntP("wait", "w", DefaultScanTime, "Duration in seconds to scan for devices")
	return findCmd
}

func runFind(devices []roku.Device) error {
	if len(devices) == 0 {
		fmt.Println("No Roku devices found on your network.")
		return nil
	}

	p := tea.NewProgram(initialFindModel(devices))
	m, err := p.Run()
	if err != nil {
		return err
	}

	finalModel := m.(findModel)
	if finalModel.selected >= 0 {
		selectedIP := devices[finalModel.selected].IP
		var deviceIPs []string
		for _, device := range devices {
			deviceIPs = append(deviceIPs, device.IP)
		}
		viper.Set("roku.devices", deviceIPs)
		viper.Set("roku.host", selectedIP)
		AddToConfigFile()
	}
	return nil
}

type findModel struct {
	devices  []roku.Device
	cursor   int
	selected int // -1 means no selection
}

func initialFindModel(devices []roku.Device) findModel {
	return findModel{
		devices:  devices,
		cursor:   0,
		selected: -1,
	}
}

func (m findModel) Init() tea.Cmd {
	return nil
}

func (m findModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m findModel) View() string {
	s := "Select a default Roku from your network:\n\n"

	for i, device := range m.devices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, device.IP)
	}

	s += "\nPress q to quit, enter to select.\n"
	return s
}

// AddToConfigFile add a key value pair to specified .roku-remote.yaml
func AddToConfigFile() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Error finding home directory: %v\n", err)
		return
	}

	path := filepath.Join(home, ".roku-remote.yaml")
	if err := viper.WriteConfigAs(path); err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}
	fmt.Printf("Updated config file: %s\n", path)
	fmt.Printf("Default Roku device set to: %s\n", viper.GetString("roku.host"))
}
