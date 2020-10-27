package cmd

import (
	"fmt"

	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find Roku Remotes on your local network.",
	Long: `Scans local network using SSDP to locate Roku 
devices that you can interface with.
	
This command uses Simple Service Discovery Protocol (or SSDP) which
provides a mechanism where by network clients, with little or no
static configuration, can discover network services.`,
	Run: func(cmd *cobra.Command, args []string) {
		wait, _ := cmd.Flags().GetInt("wait")
		devices, err := roku.Find(wait)
		if err != nil {
			panic(fmt.Sprintln("Unable to complete (find) command", err))
		}
		handleNewDevices(devices)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().IntP("wait", "w", roku.DefaultScanDuration, "The amount of time to scan for devices.")
}
