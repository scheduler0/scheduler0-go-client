package scheduler0_go_client

import "fmt"

// GetAccountExecutionCount retrieves the execution count for an account
// accountID is used both in the URL path and as the X-Account-ID header for authentication
func (c *Client) GetAccountExecutionCount(accountID string) (*AccountExecutionCountResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/accounts/%s/execution-count", accountID), nil, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountExecutionCountResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

