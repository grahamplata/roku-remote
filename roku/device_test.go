package roku

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grahamplata/roku-remote/roku/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// customTransport redirects requests from port 8060 to the test server's port
type customTransport struct {
	testServerURL string
	base          http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the request URL to use the test server
	req.URL.Host = strings.TrimPrefix(t.testServerURL, "http://")
	req.URL.Scheme = "http"
	return t.base.RoundTrip(req)
}

func TestNewDevice(t *testing.T) {
	device := NewDevice("192.168.1.100")

	assert.NotNil(t, device)
	assert.Equal(t, "192.168.1.100", device.IP)
	assert.NotNil(t, device.Client)
}

func TestDevice_Info(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
	<specVersion><major>1</major><minor>0</minor></specVersion>
	<device>
		<friendlyName>Test Roku</friendlyName>
	</device>
</root>`))
	}))
	defer server.Close()

	device := createTestDevice(server)

	info, err := device.Info(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "Test Roku", info.Device.FriendlyName)
}

func TestDevice_DeviceInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<device-info>
	<serial-number>TEST123</serial-number>
	<model-name>Roku Ultra</model-name>
</device-info>`))
	}))
	defer server.Close()

	device := createTestDevice(server)

	deviceInfo, err := device.DeviceInfo(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "TEST123", deviceInfo.SerialNumber)
	assert.Equal(t, "Roku Ultra", deviceInfo.ModelName)
}

func TestDevice_Action(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.Path, "/keypress/")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	device := createTestDevice(server)

	err := device.Action(context.Background(), "home")

	assert.NoError(t, err)
}

func TestDevice_Launch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/launch", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	device := createTestDevice(server)

	err := device.Launch(context.Background(), "12")

	assert.NoError(t, err)
}

func TestDevice_Player(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<player state="play">
	<plugin id="12" name="Netflix"/>
</player>`))
	}))
	defer server.Close()

	device := createTestDevice(server)

	player, err := device.Player(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "play", player.State)
}

func TestDevice_Describe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<device-info>
	<serial-number>XYZ789</serial-number>
</device-info>`))
	}))
	defer server.Close()

	device := createTestDevice(server)

	deviceInfo, err := device.Describe(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "XYZ789", deviceInfo.SerialNumber)
}

func TestDevice_Install(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/install", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	device := createTestDevice(server)

	err := device.Install(context.Background(), "12")

	assert.NoError(t, err)
}

func TestDevice_FetchInstalledApps(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<apps>
	<app id="12">Netflix</app>
	<app id="13">YouTube</app>
</apps>`))
	}))
	defer server.Close()

	device := createTestDevice(server)

	apps, err := device.FetchInstalledApps(context.Background())

	require.NoError(t, err)
	assert.Len(t, apps.Apps, 2)
	assert.Equal(t, "Netflix", apps.Apps[0].Name)
}

func TestDevice_ActiveApp(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<active-app>
	<app id="12">Netflix</app>
</active-app>`))
	}))
	defer server.Close()

	device := createTestDevice(server)

	activeApp, err := device.ActiveApp(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "12", activeApp.App.ID)
	assert.Equal(t, "Netflix", activeApp.App.Name)
}

// Helper function to create a test device with a custom HTTP client
func createTestDevice(server *httptest.Server) *Device {
	httpClient := &http.Client{
		Timeout: api.DefaultTimeout,
		Transport: &customTransport{
			testServerURL: server.URL,
			base:          http.DefaultTransport,
		},
	}

	return &Device{
		IP:     "127.0.0.1",
		Client: api.NewClient("127.0.0.1", httpClient),
	}
}
