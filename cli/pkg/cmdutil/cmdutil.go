package cmdutil

import (
	"fmt"
	"net"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Helper struct{}

func NewHelper() (*Helper, error) {
	ch := &Helper{}

	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("error finding home directory: %w", err)
	}

	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".roku-remote")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Config file not found or readable: %v\n", err)
	}

	return ch, nil
}

// AddGroup adds a group of commands to the parent command.
func AddGroup(parent *cobra.Command, title string, children ...*cobra.Command) {
	group := &cobra.Group{ID: title, Title: title}
	parent.AddGroup(group)
	for _, child := range children {
		child.GroupID = title
		parent.AddCommand(child)
	}
}

// ValidateRokuHost checks if a Roku host is configured, valid, and optionally reachable
func (h *Helper) ValidateRokuHost() (string, error) {
	ip := viper.GetString("roku.host")
	if ip == "" {
		return "", fmt.Errorf("no Roku device configured. Run 'roku find' command first to set a default device")
	}

	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid host IP address: %s", ip)
	}

	// Test basic connectivity to Roku device on port 8060
	address := fmt.Sprintf("%s:8060", ip)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return "", fmt.Errorf("unable to connect to Roku device at %s: %w\n\nPlease ensure:\n  • The Roku device is powered on\n  • The device is connected to the same network\n  • The IP address is correct (run 'roku-remote device find' to re-scan)", ip, err)
	}
	conn.Close()

	return ip, nil
}
