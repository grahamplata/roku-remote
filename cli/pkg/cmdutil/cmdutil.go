package cmdutil

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Helper provides utility functions for command handling
type Helper struct{}

func NewHelper() (*Helper, error) {
	ch := &Helper{}

	// Load the config from the home directory
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("error finding home directory: %w", err)
	}

	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".roku-remote")

	viper.AutomaticEnv()

	// Read config if it exists
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Config file not found or readable: %v\n", err)
	}

	return ch, nil
}

// AddGroup adds a group of commands to the parent command.
func AddGroup(parent *cobra.Command, title string, children ...*cobra.Command) {
	group := &cobra.Group{ID: title, Title: title}

	// Add add the group to the parent command.
	parent.AddGroup(group)

	// Add the child commands to the group.
	for _, child := range children {
		child.GroupID = title
		parent.AddCommand(child)
	}
}

// ValidateRokuHost checks if a Roku host is configured and provides helpful messaging
func (h *Helper) ValidateRokuHost() (string, error) {
	ip := viper.GetString("roku.host")
	if ip == "" {
		return "", fmt.Errorf("no Roku device configured. Run 'roku find' command first to set a default device")
	}
	return ip, nil
}
