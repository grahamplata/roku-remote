package cmd

import (
	"fmt"

	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an action to your Roku Device.",
	Long:  "Using the following arguments send actions to your Roku device over your network.\n\n" + roku.AvialableActions(),
	Run: func(cmd *cobra.Command, args []string) {
		ip := viper.GetString("roku.host")
		if len(args) > 0 {
			if ip == "" {
				fmt.Println(roku.NoDefaultRoku)
			}
			r := roku.New(ip)
			r.Action(args[0])
			return
		}
		fmt.Println(roku.MissingAction)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
