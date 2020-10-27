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
	Short: "A glimpse into what is currently playing on the Roku.",
	Long:  `A glimpse into what is currently playing on the Roku.`,
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
