package roku

/*
Roku Docs
External Control Protocol (ECP)
https://developer.roku.com/docs/developer-program/debugging/external-control-api.md
*/

import (
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/koron/go-ssdp"
)

// Roku type encapsulates the roku device
type Roku struct {
	IP     string       `yaml:"ip"`
	Apps   *Apps        `yaml:"-"`
	Client *http.Client `yaml:"-"`
}

type service struct {
	ServiceType string `xml:"serviceType"`
	ServiceID   string `xml:"serviceId"`
	ControlURL  string `xml:"controlURL"`
	EventSubURL string `xml:"eventSubURL"`
	SCPDURL     string `xml:"SCPDURL"`
}

type serviceList struct {
	Services []service `xml:"service"`
}

type specVersion struct {
	Major int `xml:"major"`
	Minor int `xml:"minor"`
}

// Device type encapsulates the roku device
type Device struct {
	DeviceType       string      `xml:"deviceType"`
	FriendlyName     string      `xml:"friendlyName"`
	Manufacturer     string      `xml:"manufacturer"`
	ManufacturerURL  string      `xml:"manufacturerURL"`
	ModelDescription string      `xml:"modelDescription"`
	ModelName        string      `xml:"modelName"`
	ModelNumber      string      `xml:"modelNumber"`
	ModelURL         string      `xml:"modelURL"`
	SerialNumber     string      `xml:"serialNumber"`
	UDN              string      `xml:"UDN"`
	ServiceList      serviceList `xml:"serviceList"`
}

// Info type encapsulates the roku device info at the root endpoint
type Info struct {
	XMLName xml.Name    `xml:"root" json:"-"`
	Version specVersion `xml:"specVersion"`
	Device  Device      `xml:"device"`
}

// DeviceInfo type encapsulates
type DeviceInfo struct {
	XMLName                  xml.Name `xml:"device-info" json:"-"`
	Text                     string   `xml:",chardata" json:"text,omitempty"`
	Udn                      string   `xml:"udn" json:"udn,omitempty"`
	SerialNumber             string   `xml:"serial-number" json:"serial_number,omitempty"`
	DeviceID                 string   `xml:"device-id" json:"device_id,omitempty"`
	AdvertisingID            string   `xml:"advertising-id" json:"advertising_id,omitempty"`
	VendorName               string   `xml:"vendor-name" json:"vendor_name,omitempty"`
	ModelName                string   `xml:"model-name" json:"model_name,omitempty"`
	ModelNumber              string   `xml:"model-number" json:"model_number,omitempty"`
	ModelRegion              string   `xml:"model-region" json:"model_region,omitempty"`
	IsTv                     string   `xml:"is-tv" json:"is_tv,omitempty"`
	IsStick                  string   `xml:"is-stick" json:"is_stick,omitempty"`
	SupportsEthernet         string   `xml:"supports-ethernet" json:"supports_ethernet,omitempty"`
	WifiMac                  string   `xml:"wifi-mac" json:"wifi_mac,omitempty"`
	WifiDriver               string   `xml:"wifi-driver" json:"wifi_driver,omitempty"`
	EthernetMac              string   `xml:"ethernet-mac" json:"ethernet_mac,omitempty"`
	NetworkType              string   `xml:"network-type" json:"network_type,omitempty"`
	NetworkName              string   `xml:"network-name" json:"network_name,omitempty"`
	FriendlyDeviceName       string   `xml:"friendly-device-name" json:"friendly_device_name,omitempty"`
	FriendlyModelName        string   `xml:"friendly-model-name" json:"friendly_model_name,omitempty"`
	DefaultDeviceName        string   `xml:"default-device-name" json:"default_device_name,omitempty"`
	DeviceLocation           string   `xml:"user-device-location" json:"user_device_location,omitempty"`
	UserDeviceName           string   `xml:"user-device-name" json:"user_device_name,omitempty"`
	BuildNumber              string   `xml:"build-number" json:"build_number,omitempty"`
	SoftwareVersion          string   `xml:"software-version" json:"software_version,omitempty"`
	SoftwareBuild            string   `xml:"software-build" json:"software_build,omitempty"`
	SecureDevice             string   `xml:"secure-device" json:"secure_device,omitempty"`
	Language                 string   `xml:"language" json:"language,omitempty"`
	Country                  string   `xml:"country" json:"country,omitempty"`
	Locale                   string   `xml:"locale" json:"locale,omitempty"`
	TimeZoneAuto             string   `xml:"time-zone-auto" json:"time_zone_auto,omitempty"`
	TimeZone                 string   `xml:"time-zone" json:"time_zone,omitempty"`
	TimeZoneName             string   `xml:"time-zone-name" json:"time_zone_name,omitempty"`
	TimeZoneTz               string   `xml:"time-zone-tz" json:"time_zone_tz,omitempty"`
	TimeZoneOffset           string   `xml:"time-zone-offset" json:"time_zone_offset,omitempty"`
	ClockFormat              string   `xml:"clock-format" json:"clock_format,omitempty"`
	Uptime                   string   `xml:"uptime" json:"uptime,omitempty"`
	PowerMode                string   `xml:"power-mode" json:"power_mode,omitempty"`
	SupportsSuspend          string   `xml:"supports-suspend" json:"supports_suspend,omitempty"`
	SupportsFindRemote       string   `xml:"supports-find-remote" json:"supports_find_remote,omitempty"`
	FindRemoteIsPossible     string   `xml:"find-remote-is-possible" json:"find_remote_is_possible,omitempty"`
	SupportsAudioGuide       string   `xml:"supports-audio-guide" json:"supports_audio_guide,omitempty"`
	SupportsRva              string   `xml:"supports-rva" json:"supports_rva,omitempty"`
	DeveloperEnabled         string   `xml:"developer-enabled" json:"developer_enabled,omitempty"`
	KeyedDeveloperID         string   `xml:"keyed-developer-id" json:"keyed_developer_id,omitempty"`
	SearchEnabled            string   `xml:"search-enabled" json:"search_enabled,omitempty"`
	SearchChannelsEnabled    string   `xml:"search-channels-enabled" json:"search_channels_enabled,omitempty"`
	VoiceSearchEnabled       string   `xml:"voice-search-enabled" json:"voice_search_enabled,omitempty"`
	NotificationsEnabled     string   `xml:"notifications-enabled" json:"notifications_enabled,omitempty"`
	NotificationsFirstUse    string   `xml:"notifications-first-use" json:"notifications_first_use,omitempty"`
	SupportsPrivateListening string   `xml:"supports-private-listening" json:"supports_private_listening,omitempty"`
	HeadphonesConnected      string   `xml:"headphones-connected" json:"headphones_connected,omitempty"`
	SupportsEcsTextedit      string   `xml:"supports-ecs-textedit" json:"supports_ecs_textedit,omitempty"`
	SupportsEcsMicrophone    string   `xml:"supports-ecs-microphone" json:"supports_ecs_microphone,omitempty"`
	SupportsWakeOnWlan       string   `xml:"supports-wake-on-wlan" json:"supports_wake_on_wlan,omitempty"`
	HasPlayOnRoku            string   `xml:"has-play-on-roku" json:"has_play_on_roku,omitempty"`
	HasMobileScreensaver     string   `xml:"has-mobile-screensaver" json:"has_mobile_screensaver,omitempty"`
	SupportURL               string   `xml:"support-url" json:"support_url,omitempty"`
	GrandcentralVersion      string   `xml:"grandcentral-version" json:"grandcentral_version,omitempty"`
	TrcVersion               string   `xml:"trc-version" json:"trc_version,omitempty"`
	TrcChannelVersion        string   `xml:"trc-channel-version" json:"trc_channel_version,omitempty"`
	DavinciVersion           string   `xml:"davinci-version" json:"davinci_version,omitempty"`
}

type plugin struct {
	ID        string `xml:"id,attr"`
	Bandwidth string `xml:"bandwidth,attr"`
	Name      string `xml:"name,attr"`
}

type format struct {
	Audio    string `xml:"audio,attr"`
	Video    string `xml:"captions,attr"`
	Captions string `xml:"drm,attr"`
	DRM      string `xml:"video,attr"`
}

// Player type of roku media player
type Player struct {
	Error    string `xml:"error,attr"`
	State    string `xml:"state,attr"`
	Plugin   plugin `xml:"plugin"`
	Format   format `xml:"format"`
	Position string `xml:"position"`
	Live     bool   `xml:"is_live"`
}

// Apps type is a slice of App on a roku
type Apps struct {
	Apps []App `xml:"app"`
}

// App is a singular application on a roku
type App struct {
	Name    string `xml:",chardata" yaml:"name"`
	ID      string `xml:"id,attr" yaml:"id"`
	Type    string `xml:"type,attr" yaml:"type"`
	SubType string `xml:"subtype,attr" yaml:"sub_type"`
	Version string `xml:"version,attr" yaml:"version"`
}

// Info prints the available information about a Roku device
func (r *Roku) Info() (i *Info, err error) {
	endpoint := fmt.Sprintf(r.IP)
	resp, err := r.Client.Get(endpoint)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)

	err = xml.Unmarshal(buf, &i)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	return i, nil
}

// Describe prints the available information about a Roku device
func (r *Roku) Describe() (d *DeviceInfo, err error) {
	endpoint := fmt.Sprintf(r.IP + endpoints["device"])
	resp, err := r.Client.Get(endpoint)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)
	err = xml.Unmarshal(buf, &d)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	return d, nil
}

// Show prints the available information about a Roku device
func (d *DeviceInfo) Show() {
	resp := `Vendor:  %s
Model:   %s %s
Network: %s
MAC:     %s
Uptime:  %s
Version: v%s
`
	info := fmt.Sprintf(resp, d.VendorName, d.ModelName, d.ModelNumber, d.NetworkName, d.WifiMac, d.Uptime, d.SoftwareVersion)
	fmt.Println(info)
}

// Show prints the available information about a Roku device
func (d *Device) Show() {
	info := fmt.Sprintf("%s %s", d.ModelName, d.ModelNumber)
	fmt.Println(info)
}

// Details prints the available information about a Roku player
func (p *Player) Details() {
	app := fmt.Sprintf("Playing: %s\n", p.Plugin.Name)
	duration, _ := time.ParseDuration(strings.ReplaceAll(p.Position, " ", ""))
	watchTime := fmt.Sprintf("Watched: %s", duration)
	fmt.Println(app + watchTime)
}

// Details prints the available information about a Roku application
func (a *App) Details() {
	details := fmt.Sprintf("ID: %s Name: %s Version: %s", a.ID, a.Name, a.Version)
	fmt.Println(details)
}

// List prints all the available apps found
func (a *Apps) List() {
	fmt.Println("\nThe following apps are installed on this Roku\n---")
	for _, app := range a.Apps {
		app.Details()
	}
}

// Action pass user actions to the device
func (r *Roku) Action(action string) {
	if val, ok := actions[action]; ok {
		fmt.Printf("Sent %s action to Roku\n", action)
		endpoint := fmt.Sprintf(r.IP + endpoints["keypress"] + val)
		resp, err := r.Client.Post(endpoint, "application/json", nil)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	} else {
		fmt.Printf("%s is not available", action)
	}
}

// Launch an application on the Roku device
func (r *Roku) Launch(app string) {
	if val, ok := apps[app]; ok {
		fmt.Printf("Launching %s", app)
		endpoint := fmt.Sprintf(r.IP + endpoints["launch"] + strconv.Itoa(val))
		resp, err := r.Client.Post(endpoint, "application/json", nil)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	} else {
		fmt.Printf("%s is not available", app)
	}
}

// Player the current stream segment and position of the content being played, the running time of the content, audio format, and buffering.
func (r *Roku) Player() (a *Player, err error) {
	endpoint := fmt.Sprintf(r.IP + endpoints["player"])
	resp, err := r.Client.Get(endpoint)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)
	err = xml.Unmarshal(buf, &a)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	return a, nil
}

// FetchInstalledApps locates all the installed applications on the Roku
func (r *Roku) FetchInstalledApps() {
	a := new(Apps)
	resp, err := r.Client.Get(r.IP + endpoints["apps"])
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	buf, _ := io.ReadAll(resp.Body)
	err = xml.Unmarshal(buf, &a)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	r.Apps = a
}

// Display shows information about the Roku device
func (r *Roku) Display() {
	u, _ := url.Parse(r.IP)
	host, _, _ := net.SplitHostPort(u.Host)
	fmt.Println("Roku found:", host)
}

// Install adds the supplied application to the roku
func (r *Roku) Install(app string) {
	if val, ok := apps[app]; ok {
		fmt.Printf("Installing %s", app)
		endpoint := fmt.Sprintf(r.IP + endpoints["install"] + strconv.Itoa(val))
		resp, err := r.Client.Post(endpoint, "application/json", nil)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	} else {
		fmt.Printf("%s is not available", app)
	}
}

// New creates a new Roku Device
func New(ip string) *Roku {
	client := http.DefaultClient
	return &Roku{
		IP:     ip,
		Client: client,
	}
}

// Find roku devices via ssdp protocol
func Find(ScanDuration int) (devices []Roku, err error) {
	found, _ := ssdp.Search(RokuIdenifier, ScanDuration, "")
	for _, device := range found {
		location, _ := url.Parse(device.Location)
		r := New(location.String())
		devices = append(devices, *r)
	}
	return devices, nil
}
