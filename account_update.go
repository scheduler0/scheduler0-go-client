package scheduler0_go_client

import "fmt"

// UpdateAccount updates the name of an existing account.
// id is used both in the URL path and as the X-Account-ID header so the server-side
// account scope check (path account must match header account) passes.
func (c *Client) UpdateAccount(id string, body *AccountUpdateRequestBody) (*AccountResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s", id), body, id)
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
