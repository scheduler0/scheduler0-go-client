package scheduler0_go_client

// GetAccountAISettings retrieves the AI provider settings for a given account.
// accountID is passed as an override so the proxy can route to the right account.
func (c *Client) GetAccountAISettings(accountID string) (*AccountAISettingsResponse, error) {
	req, err := c.newRequest("GET", "/ai/settings", nil, accountID)
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
// The server acknowledges the write with a message rather than echoing back the saved
// settings; call GetAccountAISettings afterward to read back what was persisted.
func (c *Client) UpsertAccountAISettings(accountID string, settings *AccountAISettings) (*AccountAISettingsUpsertResponse, error) {
	req, err := c.newRequest("PUT", "/ai/settings", settings, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountAISettingsUpsertResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAIModels retrieves the per-provider approved model catalog from GET /ai/models.
// The catalog lists every model that Scheduler0 accepts for each provider.
func (c *Client) GetAIModels() (*AIModelsResponse, error) {
	req, err := c.newRequest("GET", "/ai/models", nil, "")
	if err != nil {
		return nil, err
	}

	var result AIModelsResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
