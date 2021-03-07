package cmd

import (
	"fmt"

	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// findCmd represents the find command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describes the currently selected Roku",
	Long: `Describes the currently selected Roku. The command
fetches details about the device like make, model and services.`,
	Run: func(cmd *cobra.Command, args []string) {
		ip := viper.GetString("roku.host")
		if ip == "" {
			fmt.Println("Consider running the find command first to set a default device")
			return
		}
		r := roku.New(ip)
		d, _ := r.Describe()
		fmt.Println(d)
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
