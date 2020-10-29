package cmd

import (
	"fmt"

	"github.com/grahamplata/roku-remote/roku"
	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func handleNewDevices(devices []roku.Roku) {
	if len(devices) > 0 {
		var items []string
		for _, device := range devices {
			items = append(items, device.IP)
		}
		prompt := promptui.Select{
			Label: "Select a default Roku from your network",
			Items: items,
		}
		_, value, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		AddToConfigFile("roku.host", value)
	} else {
		fmt.Println(roku.UnableToLocate)
	}
}

// AddToConfigFile add a key value pair to specified .roku-remote.yaml
func AddToConfigFile(key string, value string) {
	viper.Set(key, value)
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	path := fmt.Sprintf(home + "/.roku-remote.yaml")
	err = viper.WriteConfigAs(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Updated local configfile:", path)
}
