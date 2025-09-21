package device

import (
	"fmt"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

func LiveCmd(ch *cmdutil.Helper) *cobra.Command {
	var liveCmd = &cobra.Command{
		Use:   "live",
		Short: "Status of the Roku media player.",
		Long:  `Status and details about the current state of the Roku's media player.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			r := roku.NewDevice(ip)
			player, err := r.Player(ctx)
			if err != nil {
				fmt.Printf("Error getting player status: %v\n", err)
				return
			}
			fmt.Printf("Player state: %s\n", player.State)
			if player.Error != "" {
				fmt.Printf("Error: %s\n", player.Error)
			}
			if player.Plugin.Name != "" {
				fmt.Printf("Plugin: %s (ID: %s)\n", player.Plugin.Name, player.Plugin.ID)
			}
			if player.Position != "" {
				fmt.Printf("Position: %s\n", player.Position)
			}
			fmt.Printf("Live: %t\n", player.Live)
		},
	}

	return liveCmd
}
