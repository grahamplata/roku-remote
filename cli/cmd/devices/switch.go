package devices

import (
	"fmt"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SwitchCmd(ch *cmdutil.Helper) *cobra.Command {
	var switchCmd = &cobra.Command{
		Use:   "switch",
		Short: "Switch the default Roku device.",
		Long:  `Select a different Roku device from the stored list to set as the default.`,
		Run: func(cmd *cobra.Command, args []string) {
			devices := viper.GetStringSlice("roku.devices")
			if len(devices) == 0 {
				fmt.Println("No devices stored. Run 'roku find' to discover and store devices.")
				return
			}

			prompt := promptui.Select{
				Label: "Select a Roku device to set as default",
				Items: devices,
			}
			_, selectedIP, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt cancelled or failed: %v\n", err)
				return
			}

			viper.Set("roku.host", selectedIP)
			AddToConfigFile()
		},
	}

	return switchCmd
}
