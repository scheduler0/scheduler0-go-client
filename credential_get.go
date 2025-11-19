package scheduler0_go_client

import "fmt"

// GetCredential retrieves a single credential by ID
func (c *Client) GetCredential(id string) (*CredentialResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/credentials/%s", id), nil)
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

