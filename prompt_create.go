package scheduler0_go_client

// CreateJobFromPrompt creates job configurations from an AI prompt
// This endpoint requires credits and uses AI to generate job configurations
func (c *Client) CreateJobFromPrompt(body *PromptJobRequest) ([]PromptJobResponse, error) {
	req, err := c.newRequest("POST", "/prompt", body)
	if err != nil {
		return nil, err
	}

	var result []PromptJobResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

