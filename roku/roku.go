package roku

/*
Roku Docs
External Control Protocol (ECP)
https://developer.roku.com/docs/developer-program/debugging/external-control-api.md
*/

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

// Plugin type of roku media player
type Plugin struct {
	ID        string `xml:"id,attr"`
	Bandwidth string `xml:"bandwidth,attr"`
	Name      string `xml:"name,attr"`
}

// Format of roku media player
type Format struct {
	Audio    string `xml:"audio,attr"`
	Video    string `xml:"captions,attr"`
	Captions string `xml:"drm,attr"`
	DRM      string `xml:"video,attr"`
}

// Player type of roku media player
type Player struct {
	Error    string `xml:"error,attr"`
	State    string `xml:"state,attr"`
	Plugin   Plugin `xml:"plugin"`
	Format   Format `xml:"format"`
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

// DisplayAll prints all the available apps found
func (a *Apps) DisplayAll() {
	fmt.Println("\nThe following apps are installed on this Roku\n---")
	for _, app := range a.Apps {
		app.Details()
	}
}

// Action pass user actions to the device
func (r *Roku) Action(action string) {
	if val, ok := actions[action]; ok {
		fmt.Println(fmt.Sprintf("Sent %s action to Roku", action))
		endpoint := fmt.Sprintf(r.IP + endpoints["keypress"] + val)
		resp, err := r.Client.Post(endpoint, "application/json", nil)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	} else {
		fmt.Println(fmt.Sprintf("%s is not available", action))
	}
}

// Launch an application on the Roku device
func (r *Roku) Launch(app string) {
	if val, ok := apps[app]; ok {
		fmt.Println(fmt.Sprintf("Launching %s", app))
		endpoint := fmt.Sprintf(r.IP + endpoints["launch"] + strconv.Itoa(val))
		resp, err := r.Client.Post(endpoint, "application/json", nil)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	} else {
		fmt.Println(fmt.Sprintf("%s is not available", app))
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

	buf, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(buf, &a)
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
	buf, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(buf, &a)
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
		fmt.Println(fmt.Sprintf("Installing %s", app))
		endpoint := fmt.Sprintf(r.IP + endpoints["install"] + strconv.Itoa(val))
		resp, err := r.Client.Post(endpoint, "application/json", nil)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
	} else {
		fmt.Println(fmt.Sprintf("%s is not available", app))
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
