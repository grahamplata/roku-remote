package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

// Test helper: creates a mock server with custom handler
func newMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	server := httptest.NewServer(handler)

	// Create an HTTP client that redirects requests to our test server
	httpClient := &http.Client{
		Timeout: DefaultTimeout,
		Transport: &customTransport{
			testServerURL: server.URL,
			base:          http.DefaultTransport,
		},
	}

	// Use any valid IP since our custom transport will redirect
	client := NewClient("127.0.0.1", httpClient)
	return server, client
}

// Test helper: returns valid XML responses for different endpoints
func mockXMLResponse(endpoint string) string {
	responses := map[string]string{
		EndpointRoot: `<?xml version="1.0" encoding="UTF-8"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
	<specVersion><major>1</major><minor>0</minor></specVersion>
	<device>
		<deviceType>urn:roku-com:device:player:1-0</deviceType>
		<friendlyName>Roku Test</friendlyName>
		<manufacturer>Roku</manufacturer>
		<modelName>Roku Ultra</modelName>
	</device>
</root>`,
		EndpointApps: `<?xml version="1.0" encoding="UTF-8"?>
<apps>
	<app id="12" type="appl" version="5.1.0">Netflix</app>
	<app id="13" type="appl" version="1.0.0">YouTube</app>
</apps>`,
		EndpointDeviceInfo: `<?xml version="1.0" encoding="UTF-8"?>
<device-info>
	<udn>12345678</udn>
	<serial-number>YN00H5123456</serial-number>
	<device-id>ABC123</device-id>
	<vendor-name>Roku</vendor-name>
	<model-name>Roku Ultra</model-name>
	<uptime>1000</uptime>
</device-info>`,
		EndpointActiveApp: `<?xml version="1.0" encoding="UTF-8"?>
<active-app>
	<app id="12">Netflix</app>
</active-app>`,
		EndpointMediaPlayer: `<?xml version="1.0" encoding="UTF-8"?>
<player error="" state="play">
	<plugin id="12" bandwidth="1000" name="Netflix"/>
	<format audio="aac" video="h264"/>
	<position>12345</position>
	<is_live>false</is_live>
</player>`,
	}
	return responses[endpoint]
}

func TestNewClient(t *testing.T) {
	t.Run("WithCustomHTTPClient", func(t *testing.T) {
		customClient := &http.Client{Timeout: 5 * time.Second}
		client := NewClient("192.168.1.1", customClient)

		assert.NotNil(t, client)
		assert.Equal(t, "192.168.1.1", client.ip)
		assert.Equal(t, customClient, client.client)
	})

	t.Run("WithNilHTTPClient_UsesDefault", func(t *testing.T) {
		client := NewClient("192.168.1.2", nil)

		assert.NotNil(t, client)
		assert.Equal(t, "192.168.1.2", client.ip)
		assert.NotNil(t, client.client)
		assert.Equal(t, DefaultTimeout, client.client.Timeout)
	})
}

func TestClient_Info(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, EndpointRoot, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, mockXMLResponse(EndpointRoot))
		})
		defer server.Close()

		info, err := client.Info(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, 1, info.Version.Major)
		assert.Equal(t, 0, info.Version.Minor)
		assert.Equal(t, "Roku Test", info.Device.FriendlyName)
		assert.Equal(t, "Roku", info.Device.Manufacturer)
	})

	t.Run("HTTPError", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server Error")
		})
		defer server.Close()

		info, err := client.Info(context.Background())

		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Contains(t, err.Error(), "failed to get device info")
		assert.Contains(t, err.Error(), "unexpected status 500")
	})

	t.Run("InvalidXML", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "not valid xml")
		})
		defer server.Close()

		info, err := client.Info(context.Background())

		assert.Error(t, err)
		assert.Nil(t, info)
	})
}

func TestClient_Apps(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, EndpointApps, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, mockXMLResponse(EndpointApps))
		})
		defer server.Close()

		apps, err := client.Apps(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, apps)
		assert.Len(t, apps.Apps, 2)
		assert.Equal(t, "12", apps.Apps[0].ID)
		assert.Equal(t, "Netflix", apps.Apps[0].Name)
		assert.Equal(t, "13", apps.Apps[1].ID)
		assert.Equal(t, "YouTube", apps.Apps[1].Name)
	})

	t.Run("HTTPError", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		apps, err := client.Apps(context.Background())

		assert.Error(t, err)
		assert.Nil(t, apps)
		assert.Contains(t, err.Error(), "failed to get apps")
	})
}

func TestClient_DeviceInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, EndpointDeviceInfo, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, mockXMLResponse(EndpointDeviceInfo))
		})
		defer server.Close()

		deviceInfo, err := client.DeviceInfo(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, deviceInfo)
		assert.Equal(t, "12345678", deviceInfo.Udn)
		assert.Equal(t, "YN00H5123456", deviceInfo.SerialNumber)
		assert.Equal(t, "Roku", deviceInfo.VendorName)
		assert.Equal(t, "Roku Ultra", deviceInfo.ModelName)
		assert.Equal(t, int64(1000), deviceInfo.Uptime)
	})

	t.Run("HTTPError", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		})
		defer server.Close()

		deviceInfo, err := client.DeviceInfo(context.Background())

		assert.Error(t, err)
		assert.Nil(t, deviceInfo)
		assert.Contains(t, err.Error(), "failed to get device info")
	})
}

func TestClient_ActiveApp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, EndpointActiveApp, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, mockXMLResponse(EndpointActiveApp))
		})
		defer server.Close()

		activeApp, err := client.ActiveApp(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, activeApp)
		assert.Equal(t, "12", activeApp.App.ID)
		assert.Equal(t, "Netflix", activeApp.App.Name)
	})

	t.Run("LimitedModeError", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Device is in Limited mode")
		})
		defer server.Close()

		activeApp, err := client.ActiveApp(context.Background())

		assert.Error(t, err)
		assert.Nil(t, activeApp)
		assert.Contains(t, err.Error(), "Limited mode")
		assert.Contains(t, err.Error(), "Home button 5 times")
	})
}

func TestClient_MediaPlayer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, EndpointMediaPlayer, r.URL.Path)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, mockXMLResponse(EndpointMediaPlayer))
		})
		defer server.Close()

		player, err := client.MediaPlayer(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, player)
		assert.Equal(t, "play", player.State)
		assert.Equal(t, "12", player.Plugin.ID)
		assert.Equal(t, "Netflix", player.Plugin.Name)
		assert.Equal(t, "12345", player.Position)
		assert.False(t, player.Live)
	})

	t.Run("HTTPError", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})
		defer server.Close()

		player, err := client.MediaPlayer(context.Background())

		assert.Error(t, err)
		assert.Nil(t, player)
	})
}

func TestClient_Input(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, EndpointInput, r.URL.Path)
			w.WriteHeader(http.StatusOK)
		})
		defer server.Close()

		err := client.Input(context.Background(), "Hello World")

		assert.NoError(t, err)
	})

	t.Run("HTTPError", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})
		defer server.Close()

		err := client.Input(context.Background(), "test")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status 400")
	})
}

func TestClient_Search(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, EndpointSearch, r.URL.Path)
			w.WriteHeader(http.StatusOK)
		})
		defer server.Close()

		err := client.Search(context.Background(), "action movies")

		assert.NoError(t, err)
	})
}

func TestClient_Keypress(t *testing.T) {
	tests := []struct {
		name        string
		action      string
		shouldError bool
		errorMsg    string
	}{
		{"ValidAction_Home", "home", false, ""},
		{"ValidAction_Select", "select", false, ""},
		{"ValidAction_VolumeUp", "volumeup", false, ""},
		{"InvalidAction", "invalid_action", true, "invalid action"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldError {
				// No server needed for validation errors
				client := NewClient("192.168.1.1", nil)
				err := client.Keypress(context.Background(), tt.action)

				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.True(t, strings.HasPrefix(r.URL.Path, EndpointKeypress))
					w.WriteHeader(http.StatusOK)
				})
				defer server.Close()

				err := client.Keypress(context.Background(), tt.action)

				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_Keydown(t *testing.T) {
	tests := []struct {
		name        string
		action      string
		shouldError bool
		errorMsg    string
	}{
		{"ValidAction_Up", "up", false, ""},
		{"ValidAction_Down", "down", false, ""},
		{"InvalidAction", "bad_action", true, "invalid action"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldError {
				client := NewClient("192.168.1.1", nil)
				err := client.Keydown(context.Background(), tt.action)

				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.True(t, strings.HasPrefix(r.URL.Path, EndpointKeydown))
					w.WriteHeader(http.StatusOK)
				})
				defer server.Close()

				err := client.Keydown(context.Background(), tt.action)

				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_Launch(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, EndpointLaunch, r.URL.Path)
			w.WriteHeader(http.StatusOK)
		})
		defer server.Close()

		err := client.Launch(context.Background(), "12")

		assert.NoError(t, err)
	})

	t.Run("EmptyAppID", func(t *testing.T) {
		client := NewClient("192.168.1.1", nil)

		err := client.Launch(context.Background(), "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "appID cannot be empty")
	})

	t.Run("HTTPError", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer server.Close()

		err := client.Launch(context.Background(), "12")

		assert.Error(t, err)
	})
}

func TestClient_Install(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, EndpointInstall, r.URL.Path)
			w.WriteHeader(http.StatusOK)
		})
		defer server.Close()

		err := client.Install(context.Background(), "12")

		assert.NoError(t, err)
	})

	t.Run("EmptyAppID", func(t *testing.T) {
		client := NewClient("192.168.1.1", nil)

		err := client.Install(context.Background(), "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "appID cannot be empty")
	})

	t.Run("NoContentResponse", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		defer server.Close()

		err := client.Install(context.Background(), "12")

		assert.NoError(t, err)
	})
}

func TestClient_ContextCancellation(t *testing.T) {
	t.Run("CancelledContext", func(t *testing.T) {
		server, client := newMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		})
		defer server.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := client.Info(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to perform request")
	})
}
