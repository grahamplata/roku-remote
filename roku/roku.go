package roku

import (
	"fmt"
	"log"
	"net/url"

	"github.com/grahamplata/roku-remote/roku/api"
	"github.com/koron/go-ssdp"
)

// RokuIdentifier is the string to look for via ssdp
const RokuIdentifier = "roku:ecp"

// Find searches for Roku devices on the local network
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

// AvailableActions returns the available actions
func AvailableActions() map[string]string {
	return api.ExternalControlActions
}
