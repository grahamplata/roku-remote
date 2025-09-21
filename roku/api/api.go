package api

/*
Roku Docs
External Control Protocol (ECP)
https://developer.roku.com/docs/developer-program/debugging/external-control-api.md
*/

import (
	"encoding/xml"
)

const (
	EndpointRoot        = "/"
	EndpointApps        = "/query/apps"
	EndpointDeviceInfo  = "/query/device-info"
	EndpointActiveApp   = "/query/active-app"
	EndpointMediaPlayer = "/query/media-player"
	EndpointIcon        = "/query/icon/"
	EndpointInput       = "/input"
	EndpointSearch      = "/search"
	EndpointKeypress    = "/keypress"
	EndpointKeydown     = "/keydown"
	EndpointLaunch      = "/launch"
	EndpointInstall     = "/install"
)

// Info type encapsulates the roku device info at the root endpoint
type Info struct {
	XMLName xml.Name    `xml:"root" json:"-"`
	Version specVersion `xml:"specVersion" json:"version"`
	Device  Device      `xml:"device" json:"device"`
}

// Device type encapsulates the roku device
type Device struct {
	DeviceType       string      `xml:"deviceType" json:"device_type"`
	FriendlyName     string      `xml:"friendlyName" json:"friendly_name"`
	Manufacturer     string      `xml:"manufacturer" json:"manufacturer"`
	ManufacturerURL  string      `xml:"manufacturerURL" json:"manufacturer_url"`
	ModelDescription string      `xml:"modelDescription" json:"model_description"`
	ModelName        string      `xml:"modelName" json:"model_name"`
	ModelNumber      string      `xml:"modelNumber" json:"model_number"`
	ModelURL         string      `xml:"modelURL" json:"model_url"`
	SerialNumber     string      `xml:"serialNumber" json:"serial_number"`
	UDN              string      `xml:"UDN" json:"udn"`
	ServiceList      serviceList `xml:"serviceList" json:"service_list"`
}

// DeviceInfo type encapsulates detailed device information
type DeviceInfo struct {
	XMLName                  xml.Name `xml:"device-info" json:"-"`
	Text                     string   `xml:",chardata" json:"text,omitempty"`
	Udn                      string   `xml:"udn" json:"udn,omitempty"`
	SerialNumber             string   `xml:"serial-number" json:"serial_number"`
	DeviceID                 string   `xml:"device-id" json:"device_id,omitempty"`
	AdvertisingID            string   `xml:"advertising-id" json:"advertising_id,omitempty"`
	VendorName               string   `xml:"vendor-name" json:"vendor_name,omitempty"`
	ModelName                string   `xml:"model-name" json:"model_name,omitempty"`
	ModelNumber              string   `xml:"model-number" json:"model_number,omitempty"`
	ModelRegion              string   `xml:"model-region" json:"model_region,omitempty"`
	IsTv                     bool     `xml:"is-tv" json:"is_tv,omitempty"`
	IsStick                  bool     `xml:"is-stick" json:"is_stick,omitempty"`
	SupportsEthernet         bool     `xml:"supports-ethernet" json:"supports_ethernet,omitempty"`
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
	SecureDevice             bool     `xml:"secure-device" json:"secure_device,omitempty"`
	Language                 string   `xml:"language" json:"language,omitempty"`
	Country                  string   `xml:"country" json:"country,omitempty"`
	Locale                   string   `xml:"locale" json:"locale,omitempty"`
	TimeZoneAuto             string   `xml:"time-zone-auto" json:"time_zone_auto,omitempty"`
	TimeZone                 string   `xml:"time-zone" json:"time_zone,omitempty"`
	TimeZoneName             string   `xml:"time-zone-name" json:"time_zone_name,omitempty"`
	TimeZoneTz               string   `xml:"time-zone-tz" json:"time_zone_tz,omitempty"`
	TimeZoneOffset           int      `xml:"time-zone-offset" json:"time_zone_offset,omitempty"`
	ClockFormat              string   `xml:"clock-format" json:"clock_format,omitempty"`
	Uptime                   int64    `xml:"uptime" json:"uptime,omitempty"`
	PowerMode                string   `xml:"power-mode" json:"power_mode,omitempty"`
	SupportsSuspend          bool     `xml:"supports-suspend" json:"supports_suspend,omitempty"`
	SupportsFindRemote       bool     `xml:"supports-find-remote" json:"supports_find_remote,omitempty"`
	FindRemoteIsPossible     bool     `xml:"find-remote-is-possible" json:"find_remote_is_possible,omitempty"`
	SupportsAudioGuide       bool     `xml:"supports-audio-guide" json:"supports_audio_guide,omitempty"`
	SupportsRva              bool     `xml:"supports-rva" json:"supports_rva,omitempty"`
	DeveloperEnabled         bool     `xml:"developer-enabled" json:"developer_enabled,omitempty"`
	KeyedDeveloperID         string   `xml:"keyed-developer-id" json:"keyed_developer_id,omitempty"`
	SearchEnabled            bool     `xml:"search-enabled" json:"search_enabled,omitempty"`
	SearchChannelsEnabled    bool     `xml:"search-channels-enabled" json:"search_channels_enabled,omitempty"`
	VoiceSearchEnabled       bool     `xml:"voice-search-enabled" json:"voice_search_enabled,omitempty"`
	NotificationsEnabled     bool     `xml:"notifications-enabled" json:"notifications_enabled,omitempty"`
	NotificationsFirstUse    bool     `xml:"notifications-first-use" json:"notifications_first_use,omitempty"`
	SupportsPrivateListening bool     `xml:"supports-private-listening" json:"supports_private_listening,omitempty"`
	HeadphonesConnected      bool     `xml:"headphones-connected" json:"headphones_connected,omitempty"`
	SupportsEcsTextedit      bool     `xml:"supports-ecs-textedit" json:"supports_ecs_textedit,omitempty"`
	SupportsEcsMicrophone    bool     `xml:"supports-ecs-microphone" json:"supports_ecs_microphone,omitempty"`
	SupportsWakeOnWlan       string   `xml:"supports-wake-on-wlan" json:"supports_wake_on_wlan,omitempty"`
	HasPlayOnRoku            string   `xml:"has-play-on-roku" json:"has_play_on_roku,omitempty"`
	HasMobileScreensaver     string   `xml:"has-mobile-screensaver" json:"has_mobile_screensaver,omitempty"`
	SupportURL               string   `xml:"support-url" json:"support_url,omitempty"`
	GrandcentralVersion      string   `xml:"grandcentral-version" json:"grandcentral_version,omitempty"`
	TrcVersion               string   `xml:"trc-version" json:"trc_version,omitempty"`
	TrcChannelVersion        string   `xml:"trc-channel-version" json:"trc_channel_version,omitempty"`
	DavinciVersion           string   `xml:"davinci-version" json:"davinci_version,omitempty"`
}

// Apps type is a slice of App on a roku
type Apps struct {
	Apps []App `xml:"app" json:"apps"`
}

// App is a singular application on a roku
type App struct {
	Name    string `xml:",chardata" json:"name"`
	ID      string `xml:"id,attr" json:"id"`
	Type    string `xml:"type,attr" json:"type"`
	SubType string `xml:"subtype,attr" json:"sub_type"`
	Version string `xml:"version,attr" json:"version"`
}

// Player type of roku media player
type Player struct {
	Error    string `xml:"error,attr" json:"error"`
	State    string `xml:"state,attr" json:"state"`
	Plugin   plugin `xml:"plugin" json:"plugin"`
	Format   format `xml:"format" json:"format"`
	Position string `xml:"position" json:"position"`
	Live     bool   `xml:"is_live" json:"live"`
}

// ActiveApp represents the currently active application on the Roku device
type ActiveApp struct {
	App App `xml:"app" json:"app"`
}

type specVersion struct {
	Major int `xml:"major" json:"major"`
	Minor int `xml:"minor" json:"minor"`
}

type service struct {
	ServiceType string `xml:"serviceType" json:"service_type"`
	ServiceID   string `xml:"serviceId" json:"service_id"`
	ControlURL  string `xml:"controlURL" json:"control_url"`
	EventSubURL string `xml:"eventSubURL" json:"event_sub_url"`
	SCPDURL     string `xml:"SCPDURL" json:"scpd_url"`
}

type serviceList struct {
	Services []service `xml:"service" json:"services"`
}

type plugin struct {
	ID        string `xml:"id,attr" json:"id"`
	Bandwidth string `xml:"bandwidth,attr" json:"bandwidth"`
	Name      string `xml:"name,attr" json:"name"`
}

type format struct {
	Audio    string `xml:"audio,attr" json:"audio"`
	Video    string `xml:"video,attr" json:"video"`
	Captions string `xml:"captions,attr" json:"captions"`
	DRM      string `xml:"drm,attr" json:"drm"`
}
