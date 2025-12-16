package apps

import (
	"fmt"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

func ActiveCmd(ch *cmdutil.Helper) *cobra.Command {
	var activeCmd = &cobra.Command{
		Use:   "active",
		Short: "Show the currently active application on your Roku.",
		Long: `Show the currently active application on your Roku.

This command works even when the device is in Limited mode.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				return err
			}
			r := roku.NewDevice(ip)
			activeApp, err := r.ActiveApp(ctx)
			if err != nil {
				return fmt.Errorf("error getting active app: %w", err)
			}
			fmt.Printf("Active App: %s (ID: %s, Type: %s)\n", activeApp.App.Name, activeApp.App.ID, activeApp.App.Type)
			return nil
		},
	}
	return activeCmd
}
