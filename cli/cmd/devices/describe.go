package devices

import (
	"fmt"
	"time"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

func DescribeCmd(ch *cmdutil.Helper) *cobra.Command {
	var describeCmd = &cobra.Command{
		Use:   "describe",
		Short: "Describes the currently selected Roku",
		Long: `Describes the currently selected Roku. The command
fetches details about the device like make, model and services.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			r := roku.NewDevice(ip)
			info, err := r.Describe(ctx)
			if err != nil {
				fmt.Printf("Error describing device: %v\n", err)
				return
			}
			fmt.Printf("Vendor: %s\nModel: %s\nNetwork: %s\nMAC: %s\nUptime: %s\nVersion: %s\n",
				info.VendorName, info.ModelName, info.NetworkName, info.WifiMac, time.Duration(info.Uptime*int64(time.Second)), info.SoftwareVersion)
		},
	}

	return describeCmd
}
