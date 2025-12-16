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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				return err
			}

			appID := args[0]
			device := roku.NewDevice(ip)
			// Validate app exists
			apps, err := device.FetchInstalledApps(ctx)
			if err != nil {
				return fmt.Errorf("error fetching apps: %w", err)
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
				return fmt.Errorf("app '%s' not found. Use 'roku apps list' to see available apps", appID)
			}
			err = device.Launch(ctx, actualID)
			if err != nil {
				return fmt.Errorf("error launching app: %w", err)
			}
			fmt.Printf("App '%s' launched successfully.\n", appID)
			return nil
		},
	}

	return launchCmd
}
