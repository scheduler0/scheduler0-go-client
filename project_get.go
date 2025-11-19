package scheduler0_go_client

import "fmt"

// GetProject retrieves a single project by ID
func (c *Client) GetProject(id int64) (*ProjectResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/projects/%d", id), nil)
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

