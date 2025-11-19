package scheduler0_go_client

import "fmt"

// GetExecutor retrieves a single executor by ID
func (c *Client) GetExecutor(id string) (*ExecutorResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/executors/%s", id), nil)
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

