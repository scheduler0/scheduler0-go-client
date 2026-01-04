package scheduler0_go_client

import "fmt"

// IncreaseAccountExecutionCount increases the execution count for an account
// accountID is used both in the URL path and as the X-Account-ID header for authentication
func (c *Client) IncreaseAccountExecutionCount(accountID string, count uint64) (*AccountExecutionCountIncreaseResponse, error) {
	body := map[string]uint64{
		"count": count,
	}
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s/execution-count", accountID), body, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountExecutionCountIncreaseResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
