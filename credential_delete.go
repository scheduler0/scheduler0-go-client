package scheduler0_go_client

import "fmt"

// DeleteCredential deletes a credential by ID
func (c *Client) DeleteCredential(id string, body *CredentialDeleteRequestBody) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/credentials/%s", id), body)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

