package scheduler0_go_client

// GetAccountTokens retrieves the current token balance for an account.
func (c *Client) GetAccountTokens(accountID string) (*AccountTokensResponse, error) {
	req, err := c.newRequest("GET", "/accounts/"+accountID+"/tokens", nil, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountTokensResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddAccountTokens adds the given number of tokens to an account's balance.
func (c *Client) AddAccountTokens(accountID string, amount int64) (*AccountTokensAddResponse, error) {
	body := &AccountTokensAddRequest{Amount: amount}
	req, err := c.newRequest("PUT", "/accounts/"+accountID+"/tokens/add", body, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountTokensAddResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
