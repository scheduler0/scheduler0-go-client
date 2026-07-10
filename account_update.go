package scheduler0_go_client

import "fmt"

// UpdateAccount updates the name of an existing account.
func (c *Client) UpdateAccount(id string, body *AccountUpdateRequestBody) (*AccountResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s", id), body)
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
