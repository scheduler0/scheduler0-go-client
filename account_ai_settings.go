package scheduler0_go_client

// GetAccountAISettings retrieves the AI provider settings for a given account.
// accountID is passed as an override so the proxy can route to the right account.
func (c *Client) GetAccountAISettings(accountID string) (*AccountAISettingsResponse, error) {
	req, err := c.newRequest("GET", "/account/ai-settings", nil, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountAISettingsResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpsertAccountAISettings creates or updates the AI provider settings for a given account.
func (c *Client) UpsertAccountAISettings(accountID string, settings *AccountAISettings) (*AccountAISettingsResponse, error) {
	req, err := c.newRequest("PUT", "/account/ai-settings", settings, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountAISettingsResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
