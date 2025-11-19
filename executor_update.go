package scheduler0_go_client

import "fmt"

// UpdateExecutor updates an existing executor
func (c *Client) UpdateExecutor(id string, body *ExecutorUpdateRequestBody) (*ExecutorResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/executors/%s", id), body)
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

