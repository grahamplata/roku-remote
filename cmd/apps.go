package cmd

import (
	"fmt"

	"github.com/grahamplata/roku-remote/roku"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// appsCmd represents the apps command
var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "List the applications on your Roku.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ip := viper.GetString("roku.host")
		if ip == "" {
			fmt.Println("Consider running the find command first to set a default device")
			return
		}
		r := roku.New(ip)
		r.FetchInstalledApps()
		r.Apps.DisplayAll()
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
