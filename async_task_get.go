package scheduler0_go_client

import "fmt"

// GetAsyncTask retrieves an async task by request ID
func (c *Client) GetAsyncTask(requestID string) (*AsyncTaskResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/async-tasks/%s", requestID), nil)
	if err != nil {
		return nil, err
	}

	var result AsyncTaskResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

