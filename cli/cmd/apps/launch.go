package apps

import (
	"fmt"
	"strings"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

func LaunchCmd(ch *cmdutil.Helper) *cobra.Command {
	var launchCmd = &cobra.Command{
		Use:   "launch [app-id-or-name]",
		Short: "Launch applications on your Roku.",
		Long: `Launch applications on your Roku by providing an application ID or name.

Examples:
  roku apps launch 12       # Launch Netflix (app ID)
  roku apps launch netflix  # Launch by name
  
Use 'roku apps list' to see available applications and their IDs.`,
		Args: cobra.ExactArgs(1), // Ensure exactly one argument
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			appID := args[0]
			device := roku.NewDevice(ip)
			// Validate app exists
			apps, err := device.FetchInstalledApps(ctx)
			if err != nil {
				fmt.Printf("Error fetching apps: %v\n", err)
				return
			}
			var actualID string
			for _, app := range apps.Apps {
				if strings.EqualFold(app.ID, appID) {
					actualID = app.ID
					break
				} else if strings.EqualFold(app.Name, appID) {
					actualID = app.ID
					break
				}
			}
			if actualID == "" {
				fmt.Printf("App '%s' not found. Use 'roku apps list' to see available apps.\n", appID)
				return
			}
			err = device.Launch(ctx, actualID)
			if err != nil {
				fmt.Printf("Error launching app: %v\n", err)
			} else {
				fmt.Printf("App '%s' launched successfully.\n", appID)
			}
		},
	}

	return launchCmd
}
