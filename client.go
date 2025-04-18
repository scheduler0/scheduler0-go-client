package scheduler0_go_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	APIKey     string
	APISecret  string
	Version    string
}

func NewClient(baseURL, version, apiKey, apiSecret string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		BaseURL:    u,
		HTTPClient: &http.Client{},
		APIKey:     apiKey,
		APISecret:  apiSecret,
		Version:    version,
	}, nil
}

func (c *Client) newRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	versionPrefix := fmt.Sprintf("/api/%s/", c.Version)

	rel := &url.URL{Path: path.Join(fmt.Sprintf("%s%s", c.BaseURL.Path, versionPrefix), endpoint)}
	u := c.BaseURL.ResolveReference(rel)

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("X-API-Secret", c.APISecret)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

func (c *Client) ListCredentials() (*PaginatedCredentialsResponse, error) {
	req, err := c.newRequest("GET", "/credentials", nil)
	if err != nil {
		return nil, err
	}

	var result PaginatedCredentialsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Define the response structs

type PaginatedCredentialsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total       int          `json:"total"`
		Offset      int          `json:"offset"`
		Limit       int          `json:"limit"`
		Credentials []Credential `json:"credentials"`
	} `json:"data"`
}

type Credential struct {
	ID        int    `json:"id"`
	Archived  bool   `json:"archived"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	CreatedAt string `json:"date_created"`
}
