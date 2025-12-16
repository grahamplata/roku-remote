package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

const RokuPort = 8060
const DefaultTimeout = 10 * time.Second
const MaxRetries = 3
const InitialRetryDelay = 100 * time.Millisecond

// Error types for better error handling
type DeviceError struct {
	Op      string // Operation that failed
	IP      string // Device IP
	Message string // Error message
	Err     error  // Underlying error
}

func (e *DeviceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s failed for device %s: %s: %v", e.Op, e.IP, e.Message, e.Err)
	}
	return fmt.Sprintf("%s failed for device %s: %s", e.Op, e.IP, e.Message)
}

func (e *DeviceError) Unwrap() error {
	return e.Err
}

// Client is an HTTP client used to communicate with the Roku device
type Client struct {
	// IP address of the Roku device
	ip string `yaml:"ip"`
	// Client is an HTTP client used to communicate with the Roku device
	client *http.Client `yaml:"-"`
}

// NewClient creates a new API client with the provided HTTP client
func NewClient(ip string, client *http.Client) *Client {
	if client == nil {
		client = &http.Client{Timeout: DefaultTimeout}
	}
	return &Client{
		ip:     ip,
		client: client,
	}
}

// retryWithBackoff executes a function with exponential backoff retry logic
func (c *Client) retryWithBackoff(ctx context.Context, operation func() error) error {
	var lastErr error
	for attempt := 0; attempt < MaxRetries; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		// Don't retry on certain errors (non-transient)
		if strings.Contains(lastErr.Error(), "Limited mode") ||
			strings.Contains(lastErr.Error(), "unexpected status 404") ||
			strings.Contains(lastErr.Error(), "unexpected status 400") {
			return lastErr
		}

		// Don't sleep on last attempt
		if attempt < MaxRetries-1 {
			delay := InitialRetryDelay * time.Duration(math.Pow(2, float64(attempt)))
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return fmt.Errorf("operation failed after %d retries: %w", MaxRetries, lastErr)
}

// Info retrieves information about the Roku device
func (c *Client) Info(ctx context.Context) (*Info, error) {
	var info Info
	err := c.retryWithBackoff(ctx, func() error {
		return c.getAndDecode(ctx, EndpointRoot, &info)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}
	return &info, nil
}

// Apps retrieves the list of installed apps on the Roku device
func (c *Client) Apps(ctx context.Context) (*Apps, error) {
	var apps Apps
	err := c.retryWithBackoff(ctx, func() error {
		return c.getAndDecode(ctx, EndpointApps, &apps)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get apps: %w", err)
	}
	return &apps, nil
}

// DeviceInfo retrieves detailed device information from the Roku device
func (c *Client) DeviceInfo(ctx context.Context) (*DeviceInfo, error) {
	var deviceInfo DeviceInfo
	err := c.retryWithBackoff(ctx, func() error {
		return c.getAndDecode(ctx, EndpointDeviceInfo, &deviceInfo)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}
	return &deviceInfo, nil
}

// ActiveApp retrieves the currently active application on the Roku device
func (c *Client) ActiveApp(ctx context.Context) (*ActiveApp, error) {
	var activeApp ActiveApp
	err := c.retryWithBackoff(ctx, func() error {
		return c.getAndDecode(ctx, EndpointActiveApp, &activeApp)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get active app: %w", err)
	}
	return &activeApp, nil
}

// MediaPlayer retrieves the current media player state from the Roku device
func (c *Client) MediaPlayer(ctx context.Context) (*Player, error) {
	var player Player
	err := c.retryWithBackoff(ctx, func() error {
		return c.getAndDecode(ctx, EndpointMediaPlayer, &player)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get media player: %w", err)
	}
	return &player, nil
}

// Input sends text input to the Roku device
func (c *Client) Input(ctx context.Context, text string) error {
	return c.retryWithBackoff(ctx, func() error {
		return c.post(ctx, EndpointInput, fmt.Sprintf("text=%s", text))
	})
}

// Search performs a search on the Roku device
func (c *Client) Search(ctx context.Context, keyword string) error {
	return c.retryWithBackoff(ctx, func() error {
		return c.post(ctx, EndpointSearch, fmt.Sprintf("keyword=%s", keyword))
	})
}

// Keypress sends a keypress event to the Roku device
func (c *Client) Keypress(ctx context.Context, action string) error {
	val, ok := ExternalControlActions[action]
	if !ok {
		return fmt.Errorf("invalid action '%s' for device %s", action, c.ip)
	}
	return c.retryWithBackoff(ctx, func() error {
		return c.post(ctx, EndpointKeypress+val, "")
	})
}

// Keydown sends a keydown event to the Roku device (key held down)
func (c *Client) Keydown(ctx context.Context, action string) error {
	val, ok := ExternalControlActions[action]
	if !ok {
		return fmt.Errorf("invalid action '%s' for device %s", action, c.ip)
	}
	return c.retryWithBackoff(ctx, func() error {
		return c.post(ctx, EndpointKeydown+val, "")
	})
}

// Launch launches an application on the Roku device
func (c *Client) Launch(ctx context.Context, appID string) error {
	if appID == "" {
		return fmt.Errorf("appID cannot be empty for device %s", c.ip)
	}
	return c.retryWithBackoff(ctx, func() error {
		return c.post(ctx, EndpointLaunch, fmt.Sprintf("id=%s", appID))
	})
}

// Install installs an application on the Roku device
func (c *Client) Install(ctx context.Context, appID string) error {
	if appID == "" {
		return fmt.Errorf("appID cannot be empty for device %s", c.ip)
	}
	return c.retryWithBackoff(ctx, func() error {
		return c.post(ctx, EndpointInstall, fmt.Sprintf("id=%s", appID))
	})
}

func (c *Client) getAndDecode(ctx context.Context, endpoint string, target interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s:%d%s", c.ip, RokuPort, endpoint), nil)
	if err != nil {
		return fmt.Errorf("failed to create request for %s%s: %w", c.ip, endpoint, err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request to %s%s: %w", c.ip, endpoint, err)
	}
	defer resp.Body.Close()
	// Accept all 2xx status codes as success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		if resp.StatusCode == http.StatusForbidden && strings.Contains(bodyStr, "Limited mode") {
			return fmt.Errorf("roku device is in Limited mode - some commands are restricted. Try pressing the Home button 5 times quickly on your Roku remote to exit Limited mode, or use 'roku-remote active-app' to see the current app")
		}
		return fmt.Errorf("unexpected status %d from %s%s: %s", resp.StatusCode, c.ip, endpoint, bodyStr)
	}
	return xml.NewDecoder(resp.Body).Decode(target)
}

func (c *Client) post(ctx context.Context, endpoint string, data string) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://%s:%d%s", c.ip, RokuPort, endpoint), strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request for %s%s: %w", c.ip, endpoint, err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request to %s%s: %w", c.ip, endpoint, err)
	}
	defer resp.Body.Close()
	// Accept all 2xx status codes as success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		return fmt.Errorf("unexpected status %d from %s%s: %s", resp.StatusCode, c.ip, endpoint, bodyStr)
	}
	return nil
}
