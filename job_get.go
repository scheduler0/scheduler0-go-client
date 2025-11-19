package scheduler0_go_client

import "fmt"

// GetJob retrieves a single job by ID
// accountIDOverride is optional - if provided, overrides the client's default account ID
func (c *Client) GetJob(id string, accountIDOverride ...string) (*JobResponse, error) {
	var accountID string
	if len(accountIDOverride) > 0 {
		accountID = accountIDOverride[0]
	}
	req, err := c.newRequest("GET", fmt.Sprintf("/jobs/%s", id), nil, accountID)
	if err != nil {
		return nil, err
	}

	var result JobResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

