package scheduler0_go_client

import (
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	APIKey     string
	APISecret  string
	Version    string
	// Basic Auth for peer communication
	Username string
	Password string
	// Account ID for most endpoints
	AccountID string
}

func NewClient(baseURL, version string, options ...ClientOption) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL:    u,
		HTTPClient: &http.Client{},
		Version:    version,
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client, nil
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithAPIKey sets the API key and secret for authentication
func WithAPIKey(apiKey, apiSecret string) ClientOption {
	return func(c *Client) {
		c.APIKey = apiKey
		c.APISecret = apiSecret
	}
}

// WithAccountID sets the account ID for requests
func WithAccountID(accountID string) ClientOption {
	return func(c *Client) {
		c.AccountID = accountID
	}
}

// WithBasicAuth sets the username and password for basic authentication
func WithBasicAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.Username = username
		c.Password = password
	}
}

// Convenience functions for common use cases

// NewAPIClient creates a client with API key authentication
func NewAPIClient(baseURL, version, apiKey, apiSecret string) (*Client, error) {
	return NewClient(baseURL, version, WithAPIKey(apiKey, apiSecret))
}

// NewAPIClientWithAccount creates a client with API key authentication and account ID
func NewAPIClientWithAccount(baseURL, version, apiKey, apiSecret, accountID string) (*Client, error) {
	return NewClient(baseURL, version, WithAPIKey(apiKey, apiSecret), WithAccountID(accountID))
}

// NewBasicAuthClient creates a client with basic authentication for peer communication
func NewBasicAuthClient(baseURL, version, username, password string) (*Client, error) {
	return NewClient(baseURL, version, WithBasicAuth(username, password))
}
