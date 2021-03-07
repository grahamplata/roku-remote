package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// findCmd represents the find command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describes the currently selected Roku",
	Long: `Describes the currently selected Roku. The command
fetches details about the device like make, model and services.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("describe")
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
