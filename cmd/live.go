package cmd

import (
	"fmt"

	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// liveCmd represents the live command
var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "Stats about the devices media player.",
	Long:  `Stats and details about the current state of the Roku's media player.`,
	Run: func(cmd *cobra.Command, args []string) {
		ip := viper.GetString("roku.host")
		if ip == "" {
			fmt.Println("Consider running the find command first to set a default device")
			return
		}
		r := roku.New(ip)
		p, _ := r.Player()
		p.Details()
	},
}

func init() {
	rootCmd.AddCommand(liveCmd)
}
