package scheduler0_go_client

import "fmt"

// GetAccountPromptCount retrieves the remaining monthly AI prompt-request quota for an account.
// accountID is used both in the URL path and as the X-Account-ID header for authentication.
func (c *Client) GetAccountPromptCount(accountID string) (*AccountPromptCountResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/accounts/%s/prompt-count", accountID), nil, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountPromptCountResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IncreaseAccountPromptCount increases the remaining prompt-request quota for an account.
// accountID is used both in the URL path and as the X-Account-ID header for authentication.
func (c *Client) IncreaseAccountPromptCount(accountID string, count uint64) (*AccountPromptCountIncreaseResponse, error) {
	body := map[string]uint64{
		"count": count,
	}
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s/prompt-count", accountID), body, accountID)
	if err != nil {
		return nil, err
	}

	var result AccountPromptCountIncreaseResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
