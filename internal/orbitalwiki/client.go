// Package orbitalwiki is a client for the Orbital Wiki REST API.
package orbitalwiki

import (
	"net/http"
	"time"
)

// DefaultBaseURL is the production Orbital Wiki API.
const DefaultBaseURL = "https://www.orbitalwiki.com/api/v1"

// Client holds what's needed to call the Orbital Wiki API.
type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides DefaultBaseURL.
func WithBaseURL(u string) Option { return func(c *Client) { c.BaseURL = u } }

// WithHTTPClient replaces the default *http.Client.
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.HTTP = h } }

// NewClient returns a Client with DefaultBaseURL and an *http.Client with a 15s timeout,
// modified by opts.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
