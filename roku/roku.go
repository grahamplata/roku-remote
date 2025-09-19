package roku

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/grahamplata/roku-remote/roku/api"
	"github.com/koron/go-ssdp"
)

// RokuIdentifier is the string to look for via ssdp
const RokuIdentifier = "roku:ecp"

func Find(ScanDuration int) (devices []Device, err error) {
	found, err := ssdp.Search(RokuIdentifier, ScanDuration, "")
	if err != nil {
		return nil, fmt.Errorf("failed to search for Roku devices: %w", err)
	}
	for _, device := range found {
		location, err := url.Parse(device.Location)
		if err != nil {
			log.Printf("Failed to parse device location %s: %v", device.Location, err)
			continue
		}
		r := NewDevice(location.Hostname())
		devices = append(devices, *r)
		log.Printf("Found Roku device at %s", location.Hostname())
	}
	return devices, nil
}

// AvailableActions returns a formatted string of available actions
func AvailableActions() string {
	actions := []string{}
	for action := range api.ExternalControlActions {
		actions = append(actions, action)
	}
	return "Available actions: " + strings.Join(actions, ", ")
}
