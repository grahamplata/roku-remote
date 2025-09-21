package device

import (
	"fmt"
	"os"
	"strings"

	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
)

func SendCmd(ch *cmdutil.Helper) *cobra.Command {
	var sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Send an action to your Roku Device.",
		Long:  "Using the following arguments send actions to your Roku device over your network.\n\n" + showAvailableActions(),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ip, err := ch.ValidateRokuHost()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			if len(args) > 0 {
				input := strings.ToLower(args[0])
				available := roku.AvailableActions()
				_, ok := available[input]
				if !ok {
					fmt.Printf("Action '%s' is not recognized. %s\n", input, showAvailableActions())
					os.Exit(1)
				}
				r := roku.NewDevice(ip)
				err = r.Action(ctx, input)
				if err != nil {
					fmt.Printf("Error sending action: %v\n", err)
					return
				}
				fmt.Println("Action sent successfully.")
				return
			}
			fmt.Println("You need to include an action. Try using --help")
			os.Exit(1)
		},
	}
	return sendCmd
}

func showAvailableActions() string {
	actions := roku.AvailableActions()
	var b strings.Builder
	b.WriteString("Available Actions:\n")
	for action := range actions {
		b.WriteString(fmt.Sprintf(" - %s\n", action))
	}
	return b.String()
}
