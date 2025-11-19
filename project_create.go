package scheduler0_go_client

// CreateProject creates a new project
func (c *Client) CreateProject(body *ProjectRequestBody) (*ProjectResponse, error) {
	req, err := c.newRequest("POST", "/projects", body)
	if err != nil {
		return nil, err
	}

	var result ProjectResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

