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
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if len(args) == 0 {
				fmt.Println("You must provide an application name or id to add.")
				return
			}
			appID := strings.TrimSpace(args[0])
			if appID == "" {
				fmt.Println("You must provide an application name or id to add.")
				return
			}
			device := roku.NewDevice(ip)
			// Check if app already exists
			apps, err := device.FetchInstalledApps(ctx)
			if err != nil {
				fmt.Printf("Error fetching apps: %v\n", err)
				return
			}
			for _, app := range apps.Apps {
				if strings.EqualFold(app.ID, appID) || strings.EqualFold(app.Name, appID) {
					fmt.Printf("App '%s' is already installed.\n", appID)
					return
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
				fmt.Printf("App '%s' not found in Roku store. Please check the app name or use an app ID.\n", appID)
				return
			}
			err = device.Install(ctx, actualID)
			if err != nil {
				fmt.Printf("Error installing app: %v\n", err)
			} else {
				fmt.Println("App installed successfully.")
			}
		},
	}

	return addCmd
}
