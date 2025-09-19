package roku

import (
	"context"
	"net/http"

	"github.com/grahamplata/roku-remote/roku/api"
)

// Device represents a Roku device on the network
type Device struct {
	// IP address of the Roku device
	IP string `yaml:"ip"`
	// Client is an HTTP client used to communicate with the Roku device
	Client *api.Client
}

// NewDevice creates a new Roku Device instance
func NewDevice(ip string) *Device {
	httpClient := &http.Client{Timeout: api.DefaultTimeout}
	return &Device{
		IP:     ip,
		Client: api.NewClient(ip, httpClient),
	}
}

// Info retrieves basic information about the Roku device
func (d *Device) Info(ctx context.Context) (*api.Info, error) {
	return d.Client.Info(ctx)
}

// DeviceInfo retrieves detailed information about the Roku device
func (d *Device) DeviceInfo(ctx context.Context) (*api.DeviceInfo, error) {
	return d.Client.DeviceInfo(ctx)
}

// Action issues a remote input action to the Roku device with retries
func (d *Device) Action(ctx context.Context, action string) error {
	return d.Client.Keypress(ctx, action)
}

// Launch an application on the Roku device
func (d *Device) Launch(ctx context.Context, appID string) error {
	return d.Client.Launch(ctx, appID)
}

// Player retrieves the current media player state
func (d *Device) Player(ctx context.Context) (*api.Player, error) {
	return d.Client.MediaPlayer(ctx)
}

// Describe retrieves and formats device details
func (d *Device) Describe(ctx context.Context) (*api.DeviceInfo, error) {
	return d.Client.DeviceInfo(ctx)
}

// Install installs an application on the Roku device
func (d *Device) Install(ctx context.Context, appID string) error {
	return d.Client.Install(ctx, appID)
}

// FetchInstalledApps retrieves the list of installed apps
func (d *Device) FetchInstalledApps(ctx context.Context) (*api.Apps, error) {
	return d.Client.Apps(ctx)
}

// ActiveApp retrieves the currently active application
func (d *Device) ActiveApp(ctx context.Context) (*api.ActiveApp, error) {
	return d.Client.ActiveApp(ctx)
}
