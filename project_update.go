package scheduler0_go_client

import "fmt"

// UpdateProject updates an existing project
func (c *Client) UpdateProject(id int64, body *ProjectUpdateRequestBody) (*ProjectResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/projects/%d", id), body)
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

