package scheduler0_go_client

import "fmt"

// GetAIUsage retrieves the account's log-derived AI request usage for the current period:
// prompt and classify limits (feature-derived), the number of successful requests used, and
// the remaining allowance, plus the period boundary. It is the single source of truth for
// rendering AI credit progress and for warning thresholds. accountID is used both in the URL
// path and as the X-Account-ID header for authentication.
func (c *Client) GetAIUsage(accountID string) (*AIUsageResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/accounts/%s/ai/usage", accountID), nil, accountID)
	if err != nil {
		return nil, err
	}

	var result AIUsageResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
