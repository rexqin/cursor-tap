// Package notify sends lightweight update notifications to the API server.
package notify

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Payload is the notify request body.
type Payload struct {
	LatestID int64 `json:"latest_id"`
}

// Client posts record update notifications to the API process.
type Client struct {
	url    string
	client *http.Client
}

// NewClient creates a notify client for the given API notify URL.
func NewClient(url string) *Client {
	return &Client{
		url: url,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

// NotifyLatest sends latest_id asynchronously; failures are ignored.
func (c *Client) NotifyLatest(latestID int64) {
	if c == nil || c.url == "" {
		return
	}
	go c.notify(latestID)
}

func (c *Client) notify(latestID int64) {
	body, err := json.Marshal(Payload{LatestID: latestID})
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}
