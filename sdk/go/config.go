package onelog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ConfigClient provides polling-based config watching for Go applications
type ConfigClient struct {
	baseURL    string
	apiKey     string
	sourceSlug string
	client     *http.Client
}

// NewConfigClient creates a new config client for fetching configuration from ULAM
//
// Example:
//
//	client := onelog.NewConfigClient(
//	    "https://api.ulam.example.com",
//	    "ulam_live_xxxxx",
//	    "my-app",
//	)
//
//	config, err := client.GetConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Database URL:", config["DATABASE_URL"])
func NewConfigClient(baseURL, apiKey, sourceSlug string) *ConfigClient {
	return &ConfigClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		sourceSlug: sourceSlug,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

// GetConfig fetches current config for the source
func (c *ConfigClient) GetConfig() (map[string]string, error) {
	url := fmt.Sprintf("%s/api/config/%s", c.baseURL, c.sourceSlug)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch config: %d", resp.StatusCode)
	}

	var result struct {
		Data map[string]string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// WatchConfig polls for config changes and invokes callback on update
func (c *ConfigClient) WatchConfig(interval time.Duration, onUpdate func(map[string]string)) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial fetch
	config, err := c.GetConfig()
	if err != nil {
		return err
	}

	lastConfig := config
	onUpdate(config)

	// Poll for changes
	for range ticker.C {
		newConfig, err := c.GetConfig()
		if err != nil {
			continue // Silently retry on error
		}

		if hasConfigChanged(lastConfig, newConfig) {
			lastConfig = newConfig
			onUpdate(newConfig)
		}
	}

	return nil
}

// WatchConfigWithContext polls for config changes with context support
func (c *ConfigClient) WatchConfigWithContext(interval time.Duration, onUpdate func(map[string]string), stopCh <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial fetch
	config, err := c.GetConfig()
	if err == nil {
		onUpdate(config)
	}

	lastConfig := config

	for {
		select {
		case <-ticker.C:
			newConfig, err := c.GetConfig()
			if err != nil {
				continue
			}

			if hasConfigChanged(lastConfig, newConfig) {
				lastConfig = newConfig
				onUpdate(newConfig)
			}
		case <-stopCh:
			return
		}
	}
}

func hasConfigChanged(old, new map[string]string) bool {
	if len(old) != len(new) {
		return true
	}

	for key, oldVal := range old {
		if newVal, exists := new[key]; !exists || newVal != oldVal {
			return true
		}
	}

	return false
}

// Example usage:
//
// client := onelog.NewConfigClient("https://ulam.example.com", "api-key", "my-app")
//
// // Blocking watch
// err := client.WatchConfig(30*time.Second, func(config map[string]string) {
//     fmt.Println("Config updated:", config)
// })
//
// // Or with context
// stopCh := make(chan struct{})
// go client.WatchConfigWithContext(30*time.Second, func(config map[string]string) {
//     fmt.Println("Config updated:", config)
// }, stopCh)
//
// // To stop watching:
// close(stopCh)
