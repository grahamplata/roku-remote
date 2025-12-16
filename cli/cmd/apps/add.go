package apps

import (
	"fmt"
	"strings"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

func AddCmd(ch *cmdutil.Helper) *cobra.Command {
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add applications to your Roku.",
		Long: `Add applications by name or id to your Roku.

Usage: roku-remote apps add netflix`,
		Args: cobra.MinimumNArgs(1), // Ensure at least one argument
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				return err
			}
			appID := strings.TrimSpace(args[0])
			if appID == "" {
				return fmt.Errorf("you must provide an application name or id to add")
			}
			device := roku.NewDevice(ip)
			// Check if app already exists
			apps, err := device.FetchInstalledApps(ctx)
			if err != nil {
				return fmt.Errorf("error fetching apps: %w", err)
			}
			for _, app := range apps.Apps {
				if strings.EqualFold(app.ID, appID) || strings.EqualFold(app.Name, appID) {
					return fmt.Errorf("app '%s' is already installed", appID)
				}
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
				return fmt.Errorf("app '%s' not found in Roku store. Please check the app name or use an app ID", appID)
			}
			err = device.Install(ctx, actualID)
			if err != nil {
				return fmt.Errorf("error installing app: %w", err)
			}
			fmt.Println("App installed successfully.")
			return nil
		},
	}

	return addCmd
}
