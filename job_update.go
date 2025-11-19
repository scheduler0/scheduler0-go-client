package scheduler0_go_client

import "fmt"

// UpdateJob updates an existing job
// accountIDOverride is optional - if provided, overrides the client's default account ID
func (c *Client) UpdateJob(id string, body *JobUpdateRequestBody, accountIDOverride ...string) (*JobResponse, error) {
	var accountID string
	if len(accountIDOverride) > 0 {
		accountID = accountIDOverride[0]
	}
	req, err := c.newRequest("PUT", fmt.Sprintf("/jobs/%s", id), body, accountID)
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

