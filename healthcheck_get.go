package scheduler0_go_client

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// Healthcheck retrieves the current leader and raft stats (no authentication required)
func (c *Client) Healthcheck() (*HealthcheckResponse, error) {
	// Create a request without authentication for healthcheck
	versionPrefix := fmt.Sprintf("/api/%s/", c.Version)
	rel := &url.URL{Path: path.Join(fmt.Sprintf("%s%s", c.BaseURL.Path, versionPrefix), "healthcheck")}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var result HealthcheckResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

