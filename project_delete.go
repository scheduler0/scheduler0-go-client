package scheduler0_go_client

import "fmt"

// DeleteProject deletes a project by ID
func (c *Client) DeleteProject(id int64, body *ProjectDeleteRequestBody) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/projects/%d", id), body)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

