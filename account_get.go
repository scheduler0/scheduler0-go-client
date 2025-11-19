package scheduler0_go_client

import "fmt"

// GetAccount retrieves a single account by ID
func (c *Client) GetAccount(id string) (*AccountResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/accounts/%s", id), nil)
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

