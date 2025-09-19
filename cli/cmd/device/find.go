package device

import (
	"fmt"
	"path/filepath"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DefaultScanDuration = 5

func FindCmd(ch *cmdutil.Helper) *cobra.Command {
	var findCmd = &cobra.Command{
		Use:   "find",
		Short: "Find Roku Remotes on your local network.",
		Long: `Scans local network using SSDP to locate Roku 
devices that you can interface with.
    
This command uses Simple Service Discovery Protocol (or SSDP) which
provides a mechanism where by network clients, with little or no
static configuration, can discover network services.`,
		Run: func(cmd *cobra.Command, args []string) {
			wait, err := cmd.Flags().GetInt("wait")
			if err != nil {
				fmt.Printf("Unable to complete (find) command: %v\n", err)
				return
			}
			fmt.Printf("Scanning for Roku devices for %d seconds...\n", wait)
			devices, err := roku.Find(wait)
			if err != nil {
				fmt.Printf("Unable to complete (find) command: %v\n", err)
				return
			}
			handleNewDevices(devices)

		},
	}
	findCmd.Flags().IntP("wait", "w", DefaultScanDuration, "Duration in seconds to scan for devices")
	return findCmd
}

func handleNewDevices(devices []roku.Device) {
	if len(devices) > 0 {
		// Store all discovered device IPs in the config
		var deviceIPs []string
		for _, device := range devices {
			deviceIPs = append(deviceIPs, device.IP)
		}
		viper.Set("roku.devices", deviceIPs)

		var items []string
		for _, device := range devices {
			items = append(items, device.IP)
		}
		prompt := promptui.Select{
			Label: "Select a default Roku from your network",
			Items: items,
		}
		_, value, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt cancelled or failed: %v\n", err)
			return
		}
		viper.Set("roku.host", value)
		AddToConfigFile()
	} else {
		fmt.Println("No Roku devices found on your network.")
	}
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
