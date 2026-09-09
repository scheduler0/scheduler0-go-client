package scheduler0_go_client

import "fmt"

// GetAccount retrieves a single account by ID.
// id is used both in the URL path and as the X-Account-ID header so the server-side
// account scope check (path account must match header account) passes.
func (c *Client) GetAccount(id string) (*AccountResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/accounts/%s", id), nil, id)
	if err != nil {
		return nil, err
	}

	var result AccountResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

