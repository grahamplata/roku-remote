package apps

import (
	"fmt"
	"sort"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

func ListCmd(ch *cmdutil.Helper) *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List the applications on your Roku.",
		Long: `List the applications on your Roku.

Usage: roku-remote apps list`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			r := roku.NewDevice(ip)
			apps, err := r.FetchInstalledApps(ctx)
			if err != nil {
				fmt.Printf("Error fetching apps: %v\n", err)
				return
			}
			// Sort apps by name
			sort.Slice(apps.Apps, func(i, j int) bool {
				return apps.Apps[i].Name < apps.Apps[j].Name
			})
			for _, app := range apps.Apps {
				fmt.Printf("%s (ID: %s)\n", app.Name, app.ID)
			}
		},
	}
	return listCmd
}
