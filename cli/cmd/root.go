package cli

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/grahamplata/roku-remote/cli/cmd/apps"
	"github.com/grahamplata/roku-remote/cli/cmd/device"
	"github.com/grahamplata/roku-remote/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.EnableCommandSorting = false
}

// Config holds CLI configuration
type Config struct {
	CfgFile string
}

func Run(ctx context.Context, version string) {
	cfg := &Config{}
	ch, err := cmdutil.NewHelper()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = RootCmd(ch, cfg, version).ExecuteContext(ctx)
	code := HandleExecuteError(ch, err)
	os.Exit(code)
}

func RootCmd(ch *cmdutil.Helper, cfg *Config, version string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "roku",
		Short:   "A cli tool to interact with roku devices on your local network.",
		Long:    `Using SSDP (Simple Service Discovery Protocol) access your Roku's RESTful API`,
		Version: version,
	}

	rootCmd.PersistentFlags().StringVar(&cfg.CfgFile, "config", "", "config file (default is $HOME/.roku-remote.yaml)")
	rootCmd.PersistentFlags().String("host", "", "host ip of the roku")
	err := viper.BindPFlag("roku.host", rootCmd.PersistentFlags().Lookup("host"))
	if err != nil {
		log.Printf("Error binding flags: %v", err)
		os.Exit(1)
	}

	// Command Groups
	// App Commands
	cmdutil.AddGroup(rootCmd, "app",
		apps.ActiveCmd(ch),
		apps.AddCmd(ch),
		apps.LaunchCmd(ch),
		apps.ListCmd(ch),
	)

	// Device Commands
	cmdutil.AddGroup(rootCmd, "device",
		device.ControlCmd(ch),
		device.DescribeCmd(ch),
		device.FindCmd(ch),
		device.LiveCmd(ch),
		device.SendCmd(ch),
		device.SwitchCmd(ch),
	)

	return rootCmd
}

func HandleExecuteError(ch *cmdutil.Helper, err error) int {
	if err == nil {
		return 0
	}

	logger := log.New(os.Stderr, "", 0)
	if strings.Contains(err.Error(), "no Roku device configured") {
		logger.Printf("Error: %v", err)
		logger.Println("Hint: Run 'roku find' to discover and configure a Roku device")
		return 1
	}

	logger.Printf("Error: %v", err)
	return 1
}
