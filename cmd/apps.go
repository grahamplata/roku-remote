package cmd

import (
	"fmt"

	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Interact with Channels on your Roku.",
	Long: `apps is for interacting with channels on your Roku

Add, Remove and List available channels.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Apps is for interacting with channels on your Roku Add, Remove and List available channels.`)
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add applications to your Roku.",
	Long: `Add applications by name or id to your Roku.

Usage: roku-remote apps add netflix`,
	Run: func(cmd *cobra.Command, args []string) {
		ip := viper.GetString("roku.host")
		if ip == "" {
			fmt.Println(roku.NoDefaultRoku)
			return
		}
		r := roku.New(ip)
		r.Install(args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the applications on your Roku.",
	Long: `List the applications on your Roku.

Usage: roku-remote apps list`,
	Run: func(cmd *cobra.Command, args []string) {
		ip := viper.GetString("roku.host")
		if ip == "" {
			fmt.Println(roku.NoDefaultRoku)
			return
		}
		r := roku.New(ip)
		r.FetchInstalledApps()
		r.Apps.DisplayAll()
	},
}

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch applications on your Roku.",
	Long: `Launch applications on your Roku.

Usage: roku-remote apps launch netflix`,
	Run: func(cmd *cobra.Command, args []string) {
		ip := viper.GetString("roku.host")
		if ip == "" {
			fmt.Println(roku.NoDefaultRoku)
			return
		}
		r := roku.New(ip)
		r.Launch(args[0])
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.AddCommand(addCmd, listCmd, launchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
