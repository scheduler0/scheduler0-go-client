package scheduler0_go_client

import "fmt"

// DeleteJob deletes a job by ID
// accountIDOverride is optional - if provided, overrides the client's default account ID
func (c *Client) DeleteJob(id string, body *JobDeleteRequestBody, accountIDOverride ...string) error {
	var accountID string
	if len(accountIDOverride) > 0 {
		accountID = accountIDOverride[0]
	}
	req, err := c.newRequest("DELETE", fmt.Sprintf("/jobs/%s", id), body, accountID)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

