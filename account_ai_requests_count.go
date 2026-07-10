package scheduler0_go_client

import "fmt"

// AccountAIRequestsCount represents the remaining monthly AI request quota for an account
type AccountAIRequestsCount struct {
	ID            int64  `json:"id"`
	AccountID     int64  `json:"accountId"`
	RequestCount  int64  `json:"requestCount"`
	DateCreated   string `json:"dateCreated"`
	DateModified  string `json:"dateModified"`
	NextResetDate string `json:"nextResetDate"`
}

// AccountAIRequestsCountResponse represents the response for account AI requests count
type AccountAIRequestsCountResponse struct {
	Success bool                   `json:"success"`
	Data    AccountAIRequestsCount `json:"data"`
}

// AccountAIRequestsCountIncreaseResponse represents the response for increasing account AI requests count
type AccountAIRequestsCountIncreaseResponse struct {
	Success bool `json:"success"`
	Data    struct {
		NewAIRequestsCount uint64 `json:"newAIRequestsCount"`
	} `json:"data"`
}

// GetAccountAIRequestsCount retrieves the remaining monthly AI request quota for an account
// accountID is used both in the URL path and as the X-Account-ID header for authentication
func (c *Client) GetAccountAIRequestsCount(accountID string) (*AccountAIRequestsCountResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/accounts/%s/ai-requests-count", accountID), nil, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountAIRequestsCountResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IncreaseAccountAIRequestsCount increases the remaining monthly AI request quota for an account
// accountID is used both in the URL path and as the X-Account-ID header for authentication
func (c *Client) IncreaseAccountAIRequestsCount(accountID string, count uint64) (*AccountAIRequestsCountIncreaseResponse, error) {
	body := map[string]uint64{
		"count": count,
	}
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s/ai-requests-count", accountID), body, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountAIRequestsCountIncreaseResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
