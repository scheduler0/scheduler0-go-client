package scheduler0_go_client

// CreateAccount creates a new account
func (c *Client) CreateAccount(body *AccountCreateRequestBody) (*AccountResponse, error) {
	req, err := c.newRequest("POST", "/accounts", body)
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

