package scheduler0_go_client

// CreateExecutor creates a new executor
func (c *Client) CreateExecutor(body *ExecutorRequestBody) (*ExecutorResponse, error) {
	req, err := c.newRequest("POST", "/executors", body)
	if err != nil {
		return nil, err
	}

	var result ExecutorResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

