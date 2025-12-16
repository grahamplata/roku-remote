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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				return err
			}
			r := roku.NewDevice(ip)
			apps, err := r.FetchInstalledApps(ctx)
			if err != nil {
				return fmt.Errorf("error fetching apps: %w", err)
			}
			// Sort apps by name
			sort.Slice(apps.Apps, func(i, j int) bool {
				return apps.Apps[i].Name < apps.Apps[j].Name
			})
			for _, app := range apps.Apps {
				fmt.Printf("%s (ID: %s)\n", app.Name, app.ID)
			}
			return nil
		},
	}
	return listCmd
}
