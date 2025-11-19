package scheduler0_go_client

import "fmt"

// UpdateCredential updates an existing credential
func (c *Client) UpdateCredential(id string, body *CredentialUpdateRequestBody) (*CredentialResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/credentials/%s", id), body)
	if err != nil {
		return nil, err
	}

	var result CredentialResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

