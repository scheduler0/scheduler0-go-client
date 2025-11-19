package scheduler0_go_client

import "fmt"

// DeleteExecutor deletes an executor by ID
func (c *Client) DeleteExecutor(id string, body *ExecutorDeleteRequestBody) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/executors/%s", id), body)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

