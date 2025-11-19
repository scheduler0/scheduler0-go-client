package scheduler0_go_client

// CreateCredential creates a new credential
func (c *Client) CreateCredential(body *CredentialCreateRequestBody) (*CredentialResponse, error) {
	req, err := c.newRequest("POST", "/credentials", body)
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

